package shopee

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	config "itemService/config"
	constants "itemService/constants"
	errors "itemService/errors"
	metrics "itemService/metrics"
	"itemService/tracing"
	"net/http"
	"strconv"

	ot "github.com/opentracing/opentracing-go"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

const (
	fetchItemPrice = "shopee.FetchItemPrice"
	component      = "itemService.external"
)

// FetchItemPrice makes a call to the Shopee API to fetch an item's information, including its name, price, image, rating etc.
// It takes in an itemID and a shopID
func FetchItemPrice(ctx context.Context, config *config.Shopee, logger *zap.Logger, itemID int64, shopID int64) (*GetItemRes, error) {
	// start tracing span from context
	// ignore outgoing context
	span, _ := ot.StartSpanFromContext(ctx, fetchItemPrice)
	// add span tags
	span.SetTag(tracing.SpanKind, tracing.SpanKindClient)
	span.SetTag(tracing.Component, tracing.ComponentExternal)
	defer span.Finish()
	// time the request
	successStr := "true" // for the metric label "success"
	errorCodeStr := "0"  // for the metric label "errorCode"
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.ExternalRequestDuration.WithLabelValues(config.GetItem.Endpoint, successStr, errorCodeStr).Observe(v)
	}))
	defer func() {
		timer.ObserveDuration()
	}()

	// make external api call
	endpoint := fmt.Sprintf("%s?itemID=%d&shopID=%d", config.GetItem.Endpoint, itemID, shopID)
	span.SetTag(tracing.PeerAddress, endpoint)
	raw, err := http.Get(endpoint)
	if err != nil {
		// error occured when making get request
		logger.Error(
			constants.ErrorExternalShopeeAPICallMsg,
			zap.String("endpoint", endpoint),
			zap.Error(err),
		)
		successStr = "false"
		return nil, errors.Error{constants.ErrorExternalShopeeAPICall, constants.ErrorExternalShopeeAPICallMsg, err}
	}
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	if err != nil {
		// io error occured
		logger.Error(
			constants.ErrorExternalShopeeAPICallMsg,
			zap.String("endpoint", endpoint),
			zap.Error(err),
		)
		return nil, err
	}

	var res *GetItemRes
	err = json.Unmarshal(body, &res)
	if err != nil {
		// unmarshalling error occured
		logger.Error(
			constants.ErrorExternalShopeeAPICallMsg,
			zap.String("endpoint", endpoint),
			zap.Error(err),
		)
		successStr = "false"
		return nil, err
	}

	errorCodeStr = strconv.Itoa(res.Error)
	logger.Info(
		constants.InfoExternalAPICall,
		zap.String("endpoint", endpoint),
		zap.Any("res", res),
	)
	return res, nil
}
