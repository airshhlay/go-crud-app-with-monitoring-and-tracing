package controllers

import (
	client "gateway/client"
	"gateway/config"
	"gateway/constants"
	req "gateway/dto/request"
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
		errorCodeInt = 400
		c.JSON(200, gin.H{
			"errorCode": 400,
			"errorMsg":  "invalid_user",
		})
		return
	}

	var addFavReq req.AddFavReq
	err := c.BindJSON(&addFavReq)
	if err != nil {
		i.logger.Error(
			"error_request_binding",
			zap.Error(err),
		)
		errorCodeInt = 400
		c.JSON(200, gin.H{
			"errorCode": 400,
			"errorMsg":  "invalid_request",
		})
		return
	}
	i.logger.Info(
		"info_request",
		zap.Any("request", addFavReq),
	)

	itemID, err := strconv.ParseInt(addFavReq.ItemID, 10, 64)
	if err != nil {
		i.logger.Error(
			"error_type_conversion",
			zap.String("itemID", addFavReq.ItemID),
			zap.Error(err),
		)
		errorCodeInt = 400
		c.JSON(200, gin.H{
			"errorCode": 400,
			"errorMsg":  "invalid_request",
		})
		return
	}

	shopID, err := strconv.ParseInt(addFavReq.ShopID, 10, 64)
	if err != nil {
		i.logger.Error(
			"error_type_conversion",
			zap.String("shopID", addFavReq.ShopID),
			zap.Error(err),
		)
		errorCodeInt = 400
		c.JSON(200, gin.H{
			"errorCode": 400,
			"errorMsg":  "invalid_request",
		})
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
		errorCodeInt = 500
		c.JSON(500, "server_error")
		return
	}
	errorCodeInt = int(clientAddFavRes.ErrorCode)

	// errorCode, if any
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
		errorCodeInt = 400
		c.JSON(200, gin.H{
			"errorCode": 400,
			"errorMsg":  "invalid_user",
		})
		return
	}

	// retrieve query params
	itemID, err := strconv.ParseInt(c.Query(constants.ItemID), 10, 64)
	if err != nil {
		i.logger.Error("error_query_params", zap.Error(err))
		errorCodeInt = 400
		c.JSON(200, gin.H{
			"errorCode": 400,
			"errorMsg":  "invalid_request",
		})
		return
	}
	shopID, err := strconv.ParseInt(c.Query(constants.ShopID), 10, 64)
	if err != nil {
		errorCodeInt = 400
		c.JSON(200, gin.H{
			"errorCode": 400,
			"errorMsg":  "invalid_request",
		})
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
		errorCodeInt = 500
		c.JSON(500, "server_error")
		return
	}
	errorCodeInt = int(clientDeleteFavRes.ErrorCode)
	// errorCode, if any
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
		errorCodeInt = 400
		c.JSON(200, gin.H{
			"errorCode": 400,
			"errorMsg":  "invalid_user",
		})
		return
	}

	// retrieve query params
	page, err := strconv.Atoi(c.Query(constants.Page))
	if err != nil {
		i.logger.Error("error_query_params", zap.Error(err))
		errorCodeInt = 400
		c.JSON(200, gin.H{
			"errorCode": 400,
			"errorMsg":  "invalid_request",
		})
	}

	if page < 0 {
		errorCodeInt = 400
		c.JSON(200, gin.H{
			"errorCode": 400,
			"errorMsg":  "invalid_request",
		})
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
		errorCodeInt = 500
		c.JSON(500, "server_error")
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
		i.logger.Error("error_no_user_id")
		c.JSON(200, gin.H{
			"errorCode": 400,
			"errorMsg":  "invalid_login_request",
		})
		return 0
	}

	userID, err := strconv.ParseInt(userIDRaw.(string), 10, 64)
	if err != nil {
		i.logger.Error("error_parse_int", zap.Error(err))
		return 0
	}
	return userID
}
