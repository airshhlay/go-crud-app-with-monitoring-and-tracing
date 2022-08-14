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

	ot "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	addFavHandler     = "gateway.AddFavHandler"
	deleteFavHandler  = "gateway.DeleteFavHandler"
	getFavListHandler = "gateway.GetFavListHandler"
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
	// start tracing span from context
	span := ot.SpanFromContext(c.Request.Context())
	i.addSpanTags(span, c)
	defer span.Finish()

	var errorCodeStr string

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

	userID := i.getUserID(c, span)
	if userID == 0 {
		errorCodeStr = strconv.Itoa(constants.ErrorGetUserIDFromToken)
		return
	}

	var addFavReq req.AddFavReq
	err := c.BindJSON(&addFavReq)
	if err != nil {
		i.logger.Info(
			constants.ErrorInvalidRequestMsg,
			zap.Error(err),
		)
		errorCodeStr = strconv.Itoa(constants.ErrorInvalidRequest)
		// add the resulting error code to the span and send a standard gateway response back to the client
		SendStandardGatewayResponse(c, span, constants.ErrorInvalidRequest, constants.ErrorInvalidRequestMsg)
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
		errorCodeStr = strconv.Itoa(constants.ErrorParseInt)
		// add the resulting error code to the span and send a standard gateway response back to the client
		SendStandardGatewayResponse(c, span, constants.ErrorParseInt, constants.ErrorParseIntMsg)
		return
	}

	shopID, err := strconv.ParseInt(addFavReq.ShopID, 10, 64)
	if err != nil {
		i.logger.Error(
			constants.ErrorParseIntMsg,
			zap.String(constants.ShopID, addFavReq.ShopID),
			zap.Error(err),
		)
		errorCodeStr = strconv.Itoa(constants.ErrorParseInt)
		// add the resulting error code to the span and send a standard gateway response back to the client
		SendStandardGatewayResponse(c, span, constants.ErrorParseInt, constants.ErrorParseIntMsg)
		return
	}

	// construct the request to be made as a grpc client to item service
	clientAddFavReq := &proto.AddFavReq{
		UserID: userID,
		ItemID: itemID,
		ShopID: shopID,
	}

	// call item service
	clientAddFavRes, err := i.client.AddFav(c.Request.Context(), clientAddFavReq)
	if err != nil {
		errorCodeStr = strconv.Itoa(constants.ErrorItemserviceConnection)
		// add the resulting error code to the span and send a standard gateway response back to the client
		SendStandardGatewayResponse(c, span, constants.ErrorItemserviceConnection, constants.ErrorItemserviceConnectionMsg)
		return
	}

	// convert error code to string for metrics label
	errorCodeStr = strconv.Itoa(int(clientAddFavRes.ErrorCode))
	// add resulting errorCode to span
	AddErrorTagsToSpan(span, clientAddFavRes.ErrorCode, clientAddFavRes.ErrorMsg)
	// return response
	c.IndentedJSON(200, clientAddFavRes)
}

// DeleteFavHandler handles requests to the /item/delete/fav endpoint
func (i *ItemServiceController) DeleteFavHandler(c *gin.Context) {
	// start tracing span from context
	span := ot.SpanFromContext(c.Request.Context())
	i.addSpanTags(span, c)
	defer span.Finish()
	var errorCodeStr string

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

	userID := i.getUserID(c, span)
	if userID == 0 {
		errorCodeStr = strconv.Itoa(constants.ErrorGetUserIDFromToken)
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
		errorCodeStr = strconv.Itoa(constants.ErrorParseInt)
		// add the resulting error code to the span and send a standard gateway response back to the client
		SendStandardGatewayResponse(c, span, constants.ErrorParseInt, constants.ErrorParseIntMsg)
		return
	}
	shopID, err := strconv.ParseInt(c.Query(constants.ShopID), 10, 64)
	if err != nil {
		i.logger.Error(
			constants.ErrorParseIntMsg,
			zap.String(constants.ShopID, c.Query(constants.ShopID)),
			zap.Error(err),
		)
		errorCodeStr = strconv.Itoa(constants.ErrorParseInt)
		// add the resulting error code to the span and send a standard gateway response back to the client
		SendStandardGatewayResponse(c, span, constants.ErrorParseInt, constants.ErrorParseIntMsg)
		return
	}

	// construct the request to be made as a grpc client to item service
	clientDeleteFavReq := &proto.DeleteFavReq{
		UserID: userID,
		ItemID: itemID,
		ShopID: shopID,
	}
	// call item service
	clientDeleteFavRes, err := i.client.DeleteFav(c.Request.Context(), clientDeleteFavReq)
	if err != nil {
		errorCodeStr = strconv.Itoa(constants.ErrorItemserviceConnection)
		// add the resulting error code to the span and send a standard gateway response back to the client
		SendStandardGatewayResponse(c, span, constants.ErrorItemserviceConnection, constants.ErrorItemserviceConnectionMsg)
		return
	}

	// convert error code to string for metrics
	errorCodeStr = strconv.Itoa(int(clientDeleteFavRes.ErrorCode))
	// add the resulting error code to the span
	AddErrorTagsToSpan(span, clientDeleteFavRes.ErrorCode, clientDeleteFavRes.ErrorMsg)
	// return response
	c.IndentedJSON(200, clientDeleteFavRes)
}

// GetFavListHandler handles requests to the /item/get/list endpoint
func (i *ItemServiceController) GetFavListHandler(c *gin.Context) {
	// start tracing span from context
	span := ot.SpanFromContext(c.Request.Context())
	i.addSpanTags(span, c)
	defer span.Finish()
	var errorCodeStr string

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

	userID := i.getUserID(c, span)
	if userID == 0 {
		errorCodeStr = strconv.Itoa(constants.ErrorGetUserIDFromToken)
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
		errorCodeStr = strconv.Itoa(constants.ErrorParseInt)
		c.JSON(200, res.GatewayResponse{ErrorCode: constants.ErrorParseInt})
		return
	}

	if page < 0 {
		errorCodeStr = strconv.Itoa(constants.ErrorInvalidRequest)
		c.JSON(200, res.GatewayResponse{ErrorCode: constants.ErrorInvalidRequest})
		return
	}

	// construct the request to be made as a grpc client to item service
	clientGetFavListReq := &proto.GetFavListReq{
		UserID: userID,
		Page:   int32(page),
	}
	// call item service
	clientGetFavListRes, err := i.client.GetFavList(c.Request.Context(), clientGetFavListReq)
	if err != nil {
		errorCodeStr = strconv.Itoa(constants.ErrorItemserviceConnection)
		// add the resulting error code to the span and send a standard gateway response back to the client
		SendStandardGatewayResponse(c, span, constants.ErrorItemserviceConnection, constants.ErrorItemserviceConnectionMsg)
		return
	}

	// convert error code to string for metrics
	errorCodeStr = strconv.Itoa(int(clientGetFavListRes.ErrorCode))

	// add the resulting error code to the span
	AddErrorTagsToSpan(span, clientGetFavListRes.ErrorCode, clientGetFavListRes.ErrorMsg)
	// return response
	c.IndentedJSON(200, clientGetFavListRes)
}

// getUserID is a helper function to retrieve the userID set by the auth middleware from the context.
func (i *ItemServiceController) getUserID(c *gin.Context, span ot.Span) int64 {
	userIDRaw, exists := c.Get(constants.UserID)

	if !exists {
		i.logger.Error(constants.ErrorNoUserIDInTokenMsg)
		SendStandardGatewayResponse(c, span, constants.ErrorNoUserIDInToken, constants.ErrorNoUserIDInTokenMsg)
		return 0
	}

	userID, err := strconv.ParseInt(userIDRaw.(string), 10, 64)
	if err != nil {
		i.logger.Error(constants.ErrorParseIntMsg, zap.Error(err))
		SendStandardGatewayResponse(c, span, constants.ErrorGetUserIDFromToken, constants.ErrorGetUserIDFromTokenMsg)
		return 0
	}
	return userID
}

func (i *ItemServiceController) addSpanTags(span ot.Span, c *gin.Context) {
	// TODO; add additional tags if needed
}
