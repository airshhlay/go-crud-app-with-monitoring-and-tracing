package shopee

import (
	"encoding/json"
	"fmt"
	"io"
	config "itemService/config"
	constants "itemService/constants"
	errors "itemService/errors"
	metrics "itemService/metrics"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

func FetchItemPrice(config *config.Shopee, logger *zap.Logger, itemId int64, shopId int64) (*ShopeeGetItemRes, error) {
	// time the request
	successStr := "true" // for the metric label "success"
	errorCodeStr := "0"  // for the metric label "errorCode"
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.ExternalRequestDuration.WithLabelValues(config.GetItem.Endpoint, successStr, errorCodeStr).Observe(v)
	}))
	defer func() {
		timer.ObserveDuration()
	}()

	// TODO: add custom error messages for io error, unmarshalling error etc
	// make external api call
	endpoint := fmt.Sprintf("%s?itemId=%d&shopId=%d", config.GetItem.Endpoint, itemId, shopId)
	raw, err := http.Get(endpoint)
	if err != nil {
		// error occured when making get request
		logger.Error(
			constants.ERROR_EXTERNAL_API_CALL_MSG,
			zap.String("endpoint", endpoint),
			zap.Error(err),
		)
		successStr = "false"
		return nil, errors.Error{constants.ERROR_EXTERNAL_API_CALL, constants.ERROR_EXTERNAL_API_CALL_MSG, err}
	}
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	if err != nil {
		// io error occured
		logger.Error(
			constants.ERROR_EXTERNAL_API_CALL_MSG,
			zap.String("endpoint", endpoint),
			zap.Error(err),
		)
		return nil, err
	}

	var res *ShopeeGetItemRes
	err = json.Unmarshal(body, &res)
	if err != nil {
		// unmarshalling error occured
		logger.Error(
			constants.ERROR_EXTERNAL_API_CALL_MSG,
			zap.String("endpoint", endpoint),
			zap.Error(err),
		)
		successStr = "false"
		return nil, err
	}

	errorCodeStr = strconv.Itoa(res.Error)
	logger.Info(
		constants.INFO_EXTERNAL_API_CALL,
		zap.String("endpoint", endpoint),
		zap.Any("res", res),
	)
	return res, nil
}
