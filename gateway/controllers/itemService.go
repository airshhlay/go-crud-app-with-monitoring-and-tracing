package controllers

import (
	client "gateway/client"
	"gateway/config"
	req "gateway/dto/request"
	metrics "gateway/metrics"
	proto "gateway/proto"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	userIdStr = "userId"
	itemIdStr = "itemId"
	shopIdStr = "shopId"
	pageStr   = "page"
)

type ItemServiceController struct {
	config *config.ItemServiceConfig
	logger *zap.Logger
	client *client.ItemServiceClient
}

func NewItemServiceController(config *config.ItemServiceConfig, logger *zap.Logger, client *client.ItemServiceClient) *ItemServiceController {
	return &ItemServiceController{
		config,
		logger,
		client,
	}
}

func (i *ItemServiceController) AddFavHandler(c *gin.Context) {
	var errorCodeStr string
	var errorCodeInt int

	// observer request latency
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.RequestLatency.WithLabelValues(i.config.Label, c.Request.URL.Path, errorCodeStr)
	}))

	defer func() {
		timer.ObserveDuration()
	}()

	userId := i.getUserId(c)
	if userId == 0 {
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

	itemId, err := strconv.ParseInt(addFavReq.ItemId, 10, 64)
	if err != nil {
		i.logger.Error(
			"error_type_conversion",
			zap.String("itemId", addFavReq.ItemId),
			zap.Error(err),
		)
		errorCodeInt = 400
		c.JSON(200, gin.H{
			"errorCode": 400,
			"errorMsg":  "invalid_request",
		})
		return
	}

	shopId, err := strconv.ParseInt(addFavReq.ShopId, 10, 64)
	if err != nil {
		i.logger.Error(
			"error_type_conversion",
			zap.String("shopId", addFavReq.ShopId),
			zap.Error(err),
		)
		errorCodeInt = 400
		c.JSON(200, gin.H{
			"errorCode": 400,
			"errorMsg":  "invalid_request",
		})
		return
	}

	clientAddFavReq := &proto.AddFavReq{
		UserId: userId,
		// UserId: 1,
		ItemId: itemId,
		ShopId: shopId,
	}

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

func (i *ItemServiceController) DeleteFavHandler(c *gin.Context) {
	var errorCodeStr string
	var errorCodeInt int

	// observer request latency
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.RequestLatency.WithLabelValues(i.config.Label, c.Request.URL.Path, errorCodeStr)
	}))
	defer func() {
		timer.ObserveDuration()
	}()

	userId := i.getUserId(c)
	if userId == 0 {
		errorCodeInt = 400
		c.JSON(200, gin.H{
			"errorCode": 400,
			"errorMsg":  "invalid_user",
		})
		return
	}

	// retrieve query params
	itemId, err := strconv.ParseInt(c.Query(itemIdStr), 10, 64)
	if err != nil {
		i.logger.Error("error_query_params", zap.Error(err))
		errorCodeInt = 400
		c.JSON(200, gin.H{
			"errorCode": 400,
			"errorMsg":  "invalid_request",
		})
		return
	}
	shopId, err := strconv.ParseInt(c.Query(shopIdStr), 10, 64)
	if err != nil {
		errorCodeInt = 400
		c.JSON(200, gin.H{
			"errorCode": 400,
			"errorMsg":  "invalid_request",
		})
		return
	}

	clientDeleteFavReq := &proto.DeleteFavReq{
		UserId: userId,
		// UserId: 1,
		ItemId: itemId,
		ShopId: shopId,
	}
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

func (i *ItemServiceController) GetFavListHandler(c *gin.Context) {
	var errorCodeStr string
	var errorCodeInt int

	// observer request latency
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.RequestLatency.WithLabelValues(i.config.Label, c.Request.URL.Path, errorCodeStr)
	}))
	defer func() {
		timer.ObserveDuration()
	}()

	userId := i.getUserId(c)
	if userId == 0 {
		errorCodeInt = 400
		c.JSON(200, gin.H{
			"errorCode": 400,
			"errorMsg":  "invalid_user",
		})
		return
	}

	// retrieve query params
	page, err := strconv.Atoi(c.Query(pageStr))
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

	clientGetFavListReq := &proto.GetFavListReq{
		UserId: userId,
		// UserId: 1,
		Page: int32(page),
	}
	clientGetFavListRes, err := i.client.GetFavList(c, clientGetFavListReq)
	if err != nil {
		errorCodeInt = 500
		c.JSON(500, "server_error")
		return
	}
	errorCodeInt = int(clientGetFavListRes.ErrorCode)
	// errorCode, if any
	errorCodeStr = strconv.Itoa(errorCodeInt)
	c.IndentedJSON(200, clientGetFavListRes)
}

func (i *ItemServiceController) Ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

func (i *ItemServiceController) getUserId(c *gin.Context) int64 {
	userIdRaw, exists := c.Get(userIdStr)

	if !exists {
		i.logger.Error("error_no_user_id")
		c.JSON(200, gin.H{
			"errorCode": 400,
			"errorMsg":  "invalid_login_request",
		})
		return 0
	}

	userId, err := strconv.ParseInt(userIdRaw.(string), 10, 64)
	if err != nil {
		i.logger.Error("error_parse_int", zap.Error(err))
		return 0
	}
	return userId
}
