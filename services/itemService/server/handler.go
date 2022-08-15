package server

import (
	"database/sql"
	"fmt"
	config "itemService/config"
	constants "itemService/constants"
	db "itemService/db"
	customErr "itemService/errors"
	shopee "itemService/external/shopee"
	pb "itemService/proto"
	util "itemService/util"
	"time"

	errGroup "golang.org/x/sync/errgroup"

	"go.uber.org/zap"
)

// Handler is a helper called by Server to handle various functions.
// It implements the bulk of the business logic.
type Handler struct {
	config       *config.Config
	dbManager    *db.DatabaseManager
	redisManager *db.RedisManager
	logger       *zap.Logger
}

// AddItemToUserFavList is called by the server when a request to the AddFav grpc service method is made
func (h *Handler) AddItemToUserFavList(itemID int64, shopID int64, userID int64) (*pb.Item, error) {
	// check if item is already in user's favourite's list
	_, query, err := h.retrieveFavFromDb(userID, itemID, shopID)

	if err != nil {
		if err != sql.ErrNoRows {
			// other unexpected error occured
			h.logger.Error(
				constants.ErrorDatabaseQueryMsg,
				zap.String(constants.Query, query),
				zap.Error(err),
			)
			return nil, &customErr.Error{constants.ErrorDatabaseQuery, constants.ErrorDatabaseQueryMsg, err}
		}
	} else {
		// item is already in user's favourites
		h.logger.Info(
			constants.InfoItemInFavourites,
			zap.Int64(constants.UserID, userID),
			zap.Int64(constants.ItemID, itemID),
			zap.Int64(constants.ShopID, shopID),
		)
		// return an error
		return nil, &customErr.Error{ErrorCode: constants.ErrorItemInFavourites, ErrorMsg: constants.InfoItemInFavourites}
	}
	// item is not yet in user favourites
	// query returned no results, item is not yet in favourites
	h.logger.Debug(
		constants.InfoItemNotInFavourites,
		zap.Int64(constants.UserID, userID),
		zap.Int64(constants.ItemID, itemID),
		zap.Int64(constants.ShopID, shopID),
	)

	// checks the cache for the item, else makes an external api call to fetch the item information
	item, err := h.getItem(itemID, shopID)
	if err != nil {
		return nil, err
	}

	// add favourite into database
	err = h.addFavIntoDb(userID, itemID, shopID)
	if err != nil {
		return nil, err
	}

	return item, err
}

// GetUserFavourites is called by the server when a request to the GetFavList grpc service method is made
func (h *Handler) GetUserFavourites(userID int64, page int32) ([]*pb.Item, int32, error) {
	favourites, err := h.retrieveFavListFromDb(userID, int(page))
	if err != nil {
		return nil, 0, err
	}

	// list of items to return to the user
	var items = make([]*pb.Item, len(favourites))

	// fetch items concurrently
	g := new(errGroup.Group)

	for i, fav := range favourites {
		fav := fav
		items := items
		i := i
		g.Go(
			func() error {
				// metrics.TotalGoRoutines.Inc()
				item, err := h.getItem(fav.ItemID, fav.ShopID)
				items[i] = item
				return err
			})
	}
	// metrics.TotalGoRoutines.Dec()
	// wait
	err = g.Wait()
	if err != nil {
		return items, 0, err
	}

	// get total pages
	totalPages, err := h.getFavouritesCount(userID)
	if err != nil {
		return items, 0, err
	}

	return items, totalPages, err
}

// DeleteFavourite is called by the server when a request to the DeleteFav grpc service method is made
func (h *Handler) DeleteFavourite(userID int64, itemID int64, shopID int64) error {
	return h.removeFavFromDb(userID, itemID, shopID)
}

// getItem is a helper function to retrieve an item's information.
// It first checks if the item is in the cache, and returns it if true.
// Else, it makes an external HTTP call to fetch the item information.
func (h *Handler) getItem(itemID int64, shopID int64) (*pb.Item, error) {
	// check if item is in cache
	item, err := h.retrieveItemFromRedis(itemID, shopID)

	if err != nil {
		// error occured with redis
		return nil, err
	}

	// item is not yet in cache, fetch item from external API
	if item == nil {
		// item not in cache, fetch item information from external api
		externalRes, err := h.fetchItemInfoFromExternal(itemID, shopID)
		if err != nil {
			return nil, err
		}

		item = &pb.Item{
			ItemID: externalRes.ItemID,
			ShopID: externalRes.ShopID,
			Price:  externalRes.Price,
			Name:   externalRes.Name,
		}

		// save item in cache
		err = h.addItemToRedis(itemID, shopID, item)
		if err != nil {
			return item, err
		}
	}

	return item, err
}

// retrieveItemFromRedis is a helper function to retrieve an item's information from redis.
// An item is identified by its itemID and shopID
func (h *Handler) retrieveItemFromRedis(itemID int64, shopID int64) (*pb.Item, error) {
	var item pb.Item

	bytes, err := h.redisManager.Get(util.FormatRedisKeyForItem(itemID, shopID))
	// unexpected error occured with redis op
	if err != nil {
		return nil, err
	}

	// unmarshal bytes
	err = util.UnmarshalProto(bytes, &item)

	if err != nil {
		// error occured when unmarshalling
		h.logger.Error(
			constants.ErrorUnmarshalMsg,
			zap.Int64(constants.ItemID, itemID),
			zap.Int64(constants.ShopID, shopID),
			zap.Error(err),
		)
		return nil, &customErr.Error{constants.ErrorUnmarshal, constants.ErrorUnmarshalMsg, err}
	}

	h.logger.Info(
		constants.InfoRedisGet,
		zap.Any(constants.Item, item),
	)

	if item.ItemID == 0 || item.ShopID == 0 || item.Price == 0 || item.Name == "" {
		// item is not in redis or incomplete
		h.logger.Info("item is not in redis")
		return nil, nil
	}

	return &item, err
}

// addItemToRedis is a helper function to add an item to redis.
func (h *Handler) addItemToRedis(itemID int64, shopID int64, item *pb.Item) error {
	// marshal into bytes to store in redis
	bytes, err := util.MarshalProto(item)
	if err != nil {
		// error occured when marshalling
		h.logger.Error(
			constants.ErrorMarshalMsg,
			zap.Int64(constants.ItemID, itemID),
			zap.Int64(constants.ShopID, shopID),
			zap.Any(constants.Item, item),
			zap.Error(err),
		)
		return &customErr.Error{constants.ErrorMarshal, constants.ErrorMarshalMsg, err}
	}

	expire := time.Duration(h.config.RedisConfig.Expire) * time.Second

	err = h.redisManager.Set(util.FormatRedisKeyForItem(itemID, shopID), bytes, expire)
	if err != nil {
		h.logger.Error(
			constants.ErrorRedisSetMsg,
			zap.Int64(constants.ItemID, itemID),
			zap.Int64(constants.ShopID, shopID),
			zap.Any(constants.Item, item),
			zap.Error(err),
		)
		return &customErr.Error{constants.ErrorRedisSet, constants.ErrorRedisSetMsg, err}
	}

	return err
}

// retrieveFavFromDb is a helper function called to perform a query on the database for a user's favourited item.
// It returns the user's favourited item as well as the used query and error.
// If the user does not have the item under their favourites, retrieveFavFromDb returns nil.
func (h *Handler) retrieveFavFromDb(userID int64, itemID int64, shopID int64) (*db.Favourite, string, error) {
	var fav db.Favourite
	query := fmt.Sprintf("SELECT * FROM Favourites WHERE userID='%d' AND itemID='%d' AND shopID='%d'", userID, itemID, shopID)
	err := h.dbManager.QueryOne(query, constants.GetFavList, &fav.ID, &fav.UserID, &fav.ItemID, &fav.ShopID, &fav.TimeAdded)
	return &fav, query, err
}

func (h *Handler) addFavIntoDb(userID int64, itemID int64, shopID int64) error {
	query := fmt.Sprintf("INSERT INTO Favourites(userID, itemID, shopID) VALUES('%d','%d','%d')", userID, itemID, shopID)
	id, err := h.dbManager.InsertRow(query, constants.AddFav)
	if err != nil {
		// error occured when inserting user into database
		return &customErr.Error{constants.ErrorDatabaseInsert, constants.ErrorDatabaseInsertMsg, err}
	}
	h.logger.Info(
		constants.InfoFavouriteAdded,
		zap.Int64(constants.UserID, userID),
		zap.Int64(constants.ItemID, itemID),
		zap.Int64(constants.ShopID, shopID),
		zap.Int64(constants.ID, id),
	)

	return err
}

// removeFavFromDb is a helper function to delete a user's favourite from the database.
func (h *Handler) removeFavFromDb(userID int64, itemID int64, shopID int64) error {
	query := fmt.Sprintf("DELETE FROM Favourites WHERE userID='%d' and itemid='%d' and shopID='%d'", userID, itemID, shopID)

	rowsDeleted, err := h.dbManager.DeleteOne(query, constants.DeleteFav)

	// unexpected error occured or no rows deleted
	if err != nil || rowsDeleted != 1 {
		if err == nil {
			// error is nil but rows deleted is not 1
			err = fmt.Errorf("rowsDeleted: %d", rowsDeleted)
		}
		return &customErr.Error{constants.ErrorDatabaseDelete, constants.ErrorDatabaseDeleteMsg, err}
	}

	return err
}

// fetchItemInfoFromExternal is a helper function for making external HTTP calls to fetch the item's information.
// If the call is successful, the item is returned.
// Else, nil is returned alongside an error.
func (h *Handler) fetchItemInfoFromExternal(itemID int64, shopID int64) (*pb.Item, error) {
	res, err := shopee.FetchItemPrice(&h.config.ExternalConfig.Shopee, h.logger, itemID, shopID)

	if err != nil {
		// external api call error
		return nil, err
	}

	if res.Error != 0 {
		// shopee api returned error
		h.logger.Error(
			constants.ErrorExternalShopeeAPICallMsg,
			zap.Int64(constants.ItemID, itemID),
			zap.Int64(constants.ShopID, shopID),
			zap.Any(constants.Res, res),
		)
		return nil, &customErr.Error{constants.ErrorExternalShopeeAPICall, constants.ErrorExternalShopeeAPICallMsg, err}
	}

	return &pb.Item{
		ItemID: res.ItemData.ItemID,
		ShopID: res.ItemData.ShopID,
		Name:   res.ItemData.Name,
		Price:  res.ItemData.Price,
	}, err
}

// retrieveFavListFromDb is a helper function to retrieve all of a user's, identified by their userID, favourites.
// It returns a list of db.Favourite
func (h *Handler) retrieveFavListFromDb(userID int64, page int) ([]db.Favourite, error) {
	query := fmt.Sprintf("SELECT * FROM Favourites WHERE userID='%d' ORDER BY timeAdded desc LIMIT %d OFFSET %d", userID, h.config.MaxPerPage, h.config.MaxPerPage*page)

	// query rows
	rows, err := h.dbManager.QueryRows(query, constants.GetFavList)
	if err != nil {
		// error occured when querying
		h.logger.Error(
			constants.ErrorDatabaseQueryMsg,
			zap.Int64(constants.UserID, userID),
			zap.Int(constants.Page, page),
			zap.String(constants.Query, query),
		)
		return nil, &customErr.Error{
			constants.ErrorDatabaseQuery, constants.ErrorDatabaseQueryMsg,
			err,
		}
	}

	var favourites []db.Favourite
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var fav db.Favourite
		err := rows.Scan(&fav.ID, &fav.UserID, &fav.ItemID, &fav.ShopID, &fav.TimeAdded)
		if err != nil {
			// error occured when scanning
			h.logger.Error(
				constants.ErrorDatabaseQueryMsg,
				zap.Error(err),
			)
			return favourites, &customErr.Error{constants.ErrorDatabaseQuery, constants.ErrorDatabaseQueryMsg, err}
		}
		favourites = append(favourites, fav)
	}
	err = rows.Err()
	if err != nil {
		h.logger.Error(
			constants.ErrorDatabaseQueryMsg,
			zap.Error(err),
		)
		return favourites, &customErr.Error{constants.ErrorDatabaseQuery, constants.ErrorDatabaseQueryMsg, err}
	}
	h.logger.Info(
		constants.InfoDatabaseQueryRows,
		zap.Int64(constants.UserID, userID),
		zap.String(constants.Query, query),
		zap.Any(constants.Res, favourites),
	)
	return favourites, err
}

// getFavouritesCount is a helper function used to count the total number of favourited items a user has.
func (h *Handler) getFavouritesCount(userID int64) (int32, error) {
	query := fmt.Sprintf("SELECT count(*) FROM Favourites WHERE userID='%d'", userID)
	var count int
	// err := h.dbManager.QueryOne(query).Scan(&count)
	err := h.dbManager.QueryOne(query, constants.GetFavCount, &count)
	if err != nil {
		return 0, &customErr.Error{constants.ErrorDatabaseQuery, constants.ErrorDatabaseQueryMsg, err}
	}

	h.logger.Info(
		constants.InfoDatabaseQuery,
		zap.Int(constants.Count, count),
		zap.String(constants.Query, query),
	)

	numPages := util.CalculateNumberOfPages(count, h.config.MaxPerPage)
	return int32(numPages), err
}
