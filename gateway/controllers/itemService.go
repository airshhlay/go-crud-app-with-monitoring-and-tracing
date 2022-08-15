package controllers

import (
	client "gateway/client"
	"gateway/config"
	"gateway/constants"
	req "gateway/dto/request"
	res "gateway/dto/response"
	metrics "gateway/metrics"
	proto "gateway/proto"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ItemServiceController is called to handle incoming HTTP requests directed to the item service.
type ItemServiceController struct {
	config *config.ItemServiceConfig
	logger *zap.Logger
	client *client.ItemServiceClient
}

// NewItemServiceController returns an ItemServiceController.
func NewItemServiceController(config *config.ItemServiceConfig, logger *zap.Logger, client *client.ItemServiceClient) *ItemServiceController {
	return &ItemServiceController{
		config,
		logger,
		client,
	}
}

// AddFavHandler handles requests to the /item/add/fav endpoint
func (i *ItemServiceController) AddFavHandler(c *gin.Context) {
	var errorCodeStr string
	var errorCodeInt int

	// observe request latency
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.RequestLatency.WithLabelValues(i.config.Label, c.Request.URL.Path, errorCodeStr).Observe(v)
	}))

	defer func() {
		timer.ObserveDuration()
	}()

	// observe response size
	responseSize := prometheus.ObserverFunc(func(v float64) {
		metrics.ResponseSize.WithLabelValues(i.config.Label, c.Request.URL.Path, errorCodeStr)
	})
	defer func() {
		responseSize.Observe(float64(c.Writer.Size()))
	}()

	userID := i.getUserID(c)
	if userID == 0 {
		errorCodeInt = constants.ErrorGetUserIDFromToken
		return
	}

	var addFavReq req.AddFavReq
	err := c.BindJSON(&addFavReq)
	if err != nil {
		i.logger.Info(
			constants.ErrorInvalidRequestMsg,
			zap.Error(err),
		)
		errorCodeInt = constants.ErrorInvalidRequest
		c.JSON(200, res.GatewayResponse{ErrorCode: constants.ErrorInvalidRequest})
		return
	}
	i.logger.Info(
		constants.InfoItemServiceRequest,
		zap.Any(constants.Request, addFavReq),
	)

	itemID, err := strconv.ParseInt(addFavReq.ItemID, 10, 64)
	if err != nil {
		i.logger.Error(
			constants.ErrorParseIntMsg,
			zap.String(constants.ItemID, addFavReq.ItemID),
			zap.Error(err),
		)
		errorCodeInt = constants.ErrorParseInt
		c.JSON(200, res.GatewayResponse{ErrorCode: constants.ErrorParseInt})
		return
	}

	shopID, err := strconv.ParseInt(addFavReq.ShopID, 10, 64)
	if err != nil {
		i.logger.Error(
			constants.ErrorParseIntMsg,
			zap.String(constants.ShopID, addFavReq.ShopID),
			zap.Error(err),
		)
		errorCodeInt = constants.ErrorParseInt
		c.JSON(200, res.GatewayResponse{ErrorCode: constants.ErrorParseInt})
		return
	}

	// construct the request to be made as a grpc client to item service
	clientAddFavReq := &proto.AddFavReq{
		UserID: userID,
		ItemID: itemID,
		ShopID: shopID,
	}

	// call item service
	clientAddFavRes, err := i.client.AddFav(c, clientAddFavReq)
	if err != nil {
		errorCodeInt = constants.ErrorItemserviceConnection
		c.JSON(500, res.GatewayResponse{ErrorCode: constants.ErrorItemserviceConnection})
		return
	}
	errorCodeInt = int(clientAddFavRes.ErrorCode)

	// convert error code to string for metrics label
	errorCodeStr = strconv.Itoa(errorCodeInt)
	c.IndentedJSON(200, clientAddFavRes)
}

// DeleteFavHandler handles requests to the /item/delete/fav endpoint
func (i *ItemServiceController) DeleteFavHandler(c *gin.Context) {
	var errorCodeStr string
	var errorCodeInt int

	// observe request latency
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.RequestLatency.WithLabelValues(i.config.Label, c.Request.URL.Path, errorCodeStr).Observe(v)
	}))
	defer func() {
		timer.ObserveDuration()
	}()

	userID := i.getUserID(c)
	if userID == 0 {
		errorCodeInt = constants.ErrorGetUserIDFromToken
		return
	}

	// retrieve query params
	itemID, err := strconv.ParseInt(c.Query(constants.ItemID), 10, 64)
	if err != nil {
		i.logger.Error(
			constants.ErrorParseIntMsg,
			zap.String(constants.ItemID, c.Query(constants.ItemID)),
			zap.Error(err),
		)
		errorCodeInt = constants.ErrorParseInt
		c.JSON(200, res.GatewayResponse{ErrorCode: constants.ErrorParseInt})
		return
	}
	shopID, err := strconv.ParseInt(c.Query(constants.ShopID), 10, 64)
	if err != nil {
		i.logger.Error(
			constants.ErrorParseIntMsg,
			zap.String(constants.ShopID, c.Query(constants.ShopID)),
			zap.Error(err),
		)
		errorCodeInt = constants.ErrorParseInt
		c.JSON(200, res.GatewayResponse{ErrorCode: constants.ErrorParseInt})
		return
	}

	// construct the request to be made as a grpc client to item service
	clientDeleteFavReq := &proto.DeleteFavReq{
		UserID: userID,
		ItemID: itemID,
		ShopID: shopID,
	}
	// call item service
	clientDeleteFavRes, err := i.client.DeleteFav(c, clientDeleteFavReq)
	if err != nil {
		errorCodeInt = constants.ErrorItemserviceConnection
		c.JSON(500, res.GatewayResponse{ErrorCode: constants.ErrorItemserviceConnection})
		return
	}

	errorCodeInt = int(clientDeleteFavRes.ErrorCode)
	// convert error code to string for metrics
	errorCodeStr = strconv.Itoa(errorCodeInt)
	c.IndentedJSON(200, clientDeleteFavRes)
}

// GetFavListHandler handles requests to the /item/get/list endpoint
func (i *ItemServiceController) GetFavListHandler(c *gin.Context) {
	var errorCodeStr string
	var errorCodeInt int

	// observe request latency
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.RequestLatency.WithLabelValues(i.config.Label, c.Request.URL.Path, errorCodeStr).Observe(v)
	}))
	defer func() {
		timer.ObserveDuration()
	}()

	// observe response size
	responseSize := prometheus.ObserverFunc(func(v float64) {
		metrics.ResponseSize.WithLabelValues(i.config.Label, c.Request.URL.Path, errorCodeStr).Observe(v)
	})
	defer func() {
		responseSize.Observe(float64(c.Writer.Size()))
	}()

	userID := i.getUserID(c)
	if userID == 0 {
		errorCodeInt = constants.ErrorGetUserIDFromToken
		return
	}

	// retrieve query params
	page, err := strconv.Atoi(c.Query(constants.Page))
	if err != nil {
		i.logger.Error(
			constants.ErrorParseIntMsg,
			zap.String(constants.Page, c.Query(constants.Page)),
			zap.Error(err),
		)
		errorCodeInt = constants.ErrorParseInt
		c.JSON(200, res.GatewayResponse{ErrorCode: constants.ErrorParseInt})
		return
	}

	if page < 0 {
		errorCodeInt = constants.ErrorInvalidRequest
		c.JSON(200, res.GatewayResponse{ErrorCode: constants.ErrorInvalidRequest})
		return
	}

	// construct the request to be made as a grpc client to item service
	clientGetFavListReq := &proto.GetFavListReq{
		UserID: userID,
		Page:   int32(page),
	}
	// call item service
	clientGetFavListRes, err := i.client.GetFavList(c, clientGetFavListReq)
	if err != nil {
		errorCodeInt = constants.ErrorItemserviceConnection
		c.JSON(500, res.GatewayResponse{ErrorCode: constants.ErrorItemserviceConnection})
		return
	}
	errorCodeInt = int(clientGetFavListRes.ErrorCode)
	errorCodeStr = strconv.Itoa(errorCodeInt)
	c.IndentedJSON(200, clientGetFavListRes)
}

// getUserID is a helper function to retrieve the userID set by the auth middleware from the context.
func (i *ItemServiceController) getUserID(c *gin.Context) int64 {
	userIDRaw, exists := c.Get(constants.UserID)

	if !exists {
		i.logger.Error(constants.ErrorNoUserIDInTokenMsg)
		c.JSON(200, res.GatewayResponse{ErrorCode: constants.ErrorNoUserIDInToken})
		return 0
	}

	userID, err := strconv.ParseInt(userIDRaw.(string), 10, 64)
	if err != nil {
		i.logger.Error(constants.ErrorParseIntMsg, zap.Error(err))
		c.JSON(200, res.GatewayResponse{ErrorCode: constants.ErrorGetUserIDFromToken})
		return 0
	}
	return userID
}
