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

const (
	deleteFavFromDbStr  = "deleteFav"
	getFavListFromDbStr = "getFavList"
	addFavIntoDbStr     = "addFav"
	getFavCount         = "getFavCount"
)

// Handler is a helper called by Server to handle various functions.
// It implements the bulk of the business logic.
type Handler struct {
	config       *config.Config
	dbManager    *db.DbManager
	redisManager *db.RedisManager
	logger       *zap.Logger
}

func (h *Handler) AddItemToUserFavList(itemId int64, shopId int64, userId int64) (*pb.Item, error) {
	h.logger.Info(
		"received request to add item",
		zap.Int64("userId", userId),
		zap.Int64("itemId", itemId),
		zap.Int64("shopId", shopId),
	)
	// check if item is already in user's favourite's list
	_, query, err := h.retrieveFavFromDb(userId, itemId, shopId)

	if err != nil {
		if err != sql.ErrNoRows {
			// other unexpected error occured
			h.logger.Error(
				constants.ERROR_DATABASE_QUERY_MSG,
				zap.String("query", query),
				zap.Error(err),
			)
			return nil, &customErr.Error{constants.ERROR_DATABASE_QUERY, constants.ERROR_DATABASE_QUERY_MSG, err}
		}
	} else {
		// item is already in user's favourites
		h.logger.Info(
			constants.INFO_ITEM_IN_FAVOURITES,
			zap.Int64("userId", userId),
			zap.Int64("itemId", itemId),
			zap.Int64("shopId", shopId),
		)
		// return an error
		return nil, &customErr.Error{ErrorCode: constants.ERROR_ITEM_IN_FAVOURITES, ErrorMsg: constants.INFO_ITEM_IN_FAVOURITES}
	}
	// item is not yet in user favourites
	// query returned no results, item is not yet in favourites
	h.logger.Debug(
		constants.INFO_ITEM_NOT_IN_FAVOURITES,
		zap.Int64("userId", userId),
		zap.Int64("itemId", itemId),
		zap.Int64("shopId", shopId),
	)

	// checks the cache for the item, else makes an external api call to fetch the item information
	item, err := h.getItem(itemId, shopId)
	if err != nil {
		return nil, err
	}

	// add favourite into database
	err = h.addFavIntoDb(userId, itemId, shopId)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (h *Handler) GetUserFavourites(userId int64, page int32) ([]*pb.Item, int32, error) {
	favourites, err := h.retrieveFavListFromDb(userId, int(page))
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
				item, err := h.getItem(fav.ItemId, fav.ShopId)
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
	totalPages, err := h.getFavouritesCount(userId)
	if err != nil {
		return items, 0, err
	}

	return items, totalPages, nil
}

func (h *Handler) DeleteFavourite(userId int64, itemId int64, shopId int64) error {
	return h.removeFavFromDb(userId, itemId, shopId)
}

func (h *Handler) getItem(itemId int64, shopId int64) (*pb.Item, error) {
	// check if item is in cache
	item, err := h.retrieveItemFromRedis(itemId, shopId)

	if err != nil {
		// error occured with redis
		return nil, err
	}

	// item is not yet in cache, fetch item from external API
	if item == nil {
		// item not in cache, fetch item information from external api
		externalRes, err := h.fetchItemInfoFromExternal(itemId, shopId)
		if err != nil {
			return nil, err
		}

		item = &pb.Item{
			ItemId: externalRes.ItemId,
			ShopId: externalRes.ShopId,
			Price:  externalRes.Price,
			Name:   externalRes.Name,
		}

		// save item in cache
		err = h.addItemToRedis(itemId, shopId, item)
		if err != nil {
			return item, err
		}
	}

	return item, err
}

func (h *Handler) retrieveItemFromRedis(itemId int64, shopId int64) (*pb.Item, error) {
	var item pb.Item

	bytes, err := h.redisManager.Get(util.FormatRedisKeyForItem(itemId, shopId))
	// unexpected error occured with redis op
	if err != nil {
		return nil, err
	}

	// unmarshal bytes
	err = util.UnmarshalProto(bytes, &item)

	if err != nil {
		// error occured when unmarshalling
		h.logger.Error(
			constants.ERROR_UNMARSHAL_MSG,
			zap.Int64("itemId", itemId),
			zap.Int64("shopId", shopId),
			zap.Error(err),
		)
		return nil, &customErr.Error{constants.ERROR_UNMARSHAL, constants.ERROR_UNMARSHAL_MSG, err}
	}

	h.logger.Info(
		constants.INFO_REDIS_GET,
		zap.Any("item", item),
	)

	if item.ItemId == 0 || item.ShopId == 0 || item.Price == 0 || item.Name == "" {
		// item is not in redis or incomplete
		h.logger.Info("item is not in redis")
		return nil, nil
	}

	return &item, nil
}

func (h *Handler) addItemToRedis(itemId int64, shopId int64, item *pb.Item) error {
	// marshal into bytes to store in redis
	bytes, err := util.MarshalProto(item)
	if err != nil {
		// error occured when marshalling
		h.logger.Error(
			constants.ERROR_MARSHAL_MSG,
			zap.Int64("itemId", itemId),
			zap.Int64("shopId", shopId),
			zap.Any("item", item),
			zap.Error(err),
		)
		return &customErr.Error{constants.ERROR_MARSHAL, constants.ERROR_MARSHAL_MSG, err}
	}

	expire := time.Duration(h.config.RedisConfig.Expire) * time.Second

	err = h.redisManager.Set(util.FormatRedisKeyForItem(itemId, shopId), bytes, expire)
	if err != nil {
		h.logger.Error(
			constants.ERROR_REDIS_SET_MSG,
			zap.Int64("itemId", itemId),
			zap.Int64("shopId", shopId),
			zap.Any("item", item),
			zap.Error(err),
		)
		return &customErr.Error{constants.ERROR_REDIS_SET, constants.ERROR_REDIS_SET_MSG, err}
	}

	return nil
}

func (h *Handler) retrieveFavFromDb(userId int64, itemId int64, shopId int64) (*db.Favourite, string, error) {
	var fav db.Favourite
	query := fmt.Sprintf("SELECT * FROM Favourites WHERE userId='%d' AND itemId='%d' AND shopId='%d'", userId, itemId, shopId)
	err := h.dbManager.QueryOne(query, getFavListStr, &fav.Id, &fav.UserId, &fav.ItemId, &fav.ShopId, &fav.TimeAdded)
	// err := res.Scan(&fav.Id, &fav.UserId, &fav.ItemId, &fav.ShopId, &fav.TimeAdded)

	return &fav, query, err
}

func (h *Handler) addFavIntoDb(userId int64, itemId int64, shopId int64) error {
	query := fmt.Sprintf("INSERT INTO Favourites(userId, itemId, shopId) VALUES('%d','%d','%d')", userId, itemId, shopId)
	id, err := h.dbManager.InsertRow(query, addFavIntoDbStr)
	if err != nil {
		// error occured when inserting user into database
		return &customErr.Error{constants.ERROR_DATABASE_INSERT, constants.ERROR_DATABASE_INSERT_MSG, err}
	}
	h.logger.Info(
		"favourite added",
		zap.Int64("userId", userId),
		zap.Int64("itemId", itemId),
		zap.Int64("shopId", shopId),
		zap.Int64("id", id),
	)

	return nil
}

func (h *Handler) removeFavFromDb(userId int64, itemId int64, shopId int64) error {
	query := fmt.Sprintf("DELETE FROM Favourites WHERE userId='%d' and itemid='%d' and shopId='%d'", userId, itemId, shopId)

	rowsDeleted, err := h.dbManager.DeleteOne(query, deleteFavFromDbStr)

	// unexpected error occured or no rows deleted
	if err != nil || rowsDeleted != 1 {
		if err == nil {
			// error is nil but rows deleted is not 1
			err = fmt.Errorf("rowsDeleted: %d", rowsDeleted)
		}
		return &customErr.Error{constants.ERROR_DATABASE_DELETE, constants.ERROR_DATABASE_DELETE_MSG, err}
	}

	return nil
}

func (h *Handler) fetchItemInfoFromExternal(itemId int64, shopId int64) (*pb.Item, error) {
	res, err := shopee.FetchItemPrice(&h.config.ExternalConfig.Shopee, h.logger, itemId, shopId)

	if err != nil {
		// external api call error
		return nil, err
	}

	if res.Error != 0 {
		// shopee api returned error
		h.logger.Error(
			constants.ERROR_SHOPEE_API_CALL_MSG,
			zap.Int64("itemId", itemId),
			zap.Int64("shopId", shopId),
			zap.Any("res", res),
		)
		return nil, &customErr.Error{constants.ERROR_EXTERNAL_SHOPEE_API_CALL, constants.ERROR_EXTERNAL_API_CALL_MSG, err}
	}

	return &pb.Item{
		ItemId: res.ItemData.ItemId,
		ShopId: res.ItemData.ShopId,
		Name:   res.ItemData.Name,
		Price:  res.ItemData.Price,
	}, nil
}

func (h *Handler) retrieveFavListFromDb(userId int64, page int) ([]db.Favourite, error) {
	query := fmt.Sprintf("SELECT * FROM Favourites WHERE userId='%d' ORDER BY timeAdded desc LIMIT %d OFFSET %d", userId, h.config.MaxPerPage, h.config.MaxPerPage*page)

	// query rows
	rows, err := h.dbManager.QueryRows(query, getFavListStr)
	if err != nil {
		// error occured when querying
		h.logger.Error(
			constants.ERROR_DATABASE_QUERY_MSG,
			zap.Int64("userId", userId),
			zap.Int("page", page),
			zap.String("query", query),
		)
		return nil, &customErr.Error{
			constants.ERROR_DATABASE_QUERY, constants.ERROR_DATABASE_QUERY_MSG,
			err,
		}
	}

	var favourites []db.Favourite
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var fav db.Favourite
		err := rows.Scan(&fav.Id, &fav.UserId, &fav.ItemId, &fav.ShopId, &fav.TimeAdded)
		if err != nil {
			// error occured when scanning
			h.logger.Error(
				constants.ERROR_DATABASE_QUERY_MSG,
				zap.Error(err),
			)
			return favourites, &customErr.Error{constants.ERROR_DATABASE_QUERY, constants.ERROR_DATABASE_QUERY_MSG, err}
		}
		favourites = append(favourites, fav)
	}
	err = rows.Err()
	if err != nil {
		h.logger.Error(
			constants.ERROR_DATABASE_QUERY_MSG,
			zap.Error(err),
		)
		return favourites, &customErr.Error{constants.ERROR_DATABASE_QUERY, constants.ERROR_DATABASE_QUERY_MSG, err}
	}
	h.logger.Info(
		constants.INFO_DATABASE_QUERY_ROWS,
		zap.Int64("userId", userId),
		zap.String("query", query),
		zap.Any("result", favourites),
	)
	return favourites, nil
}

func (h *Handler) getFavouritesCount(userId int64) (int32, error) {
	query := fmt.Sprintf("SELECT count(*) FROM Favourites WHERE userId='%d'", userId)
	var count int
	// err := h.dbManager.QueryOne(query).Scan(&count)
	err := h.dbManager.QueryOne(query, getFavCount, &count)
	if err != nil {
		return 0, &customErr.Error{constants.ERROR_DATABASE_QUERY, constants.ERROR_DATABASE_QUERY_MSG, err}
	}

	h.logger.Info(
		constants.INFO_DATABASE_QUERY,
		zap.Int("count", count),
		zap.String("query", query),
	)

	numPages := util.CalculateNumberOfPages(count, h.config.MaxPerPage)
	return int32(numPages), nil
}
