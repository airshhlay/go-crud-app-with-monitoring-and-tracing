package db

import (
	"fmt"
	"itemService/config"
	constants "itemService/constants"
	errors "itemService/errors"
	metrics "itemService/metrics"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type RedisManager struct {
	client *redis.Client
	config *config.RedisConfig
	logger *zap.Logger
}

const (
	redisOpGetStr = "get"
	redisOpSetStr = "set"
	trueStr       = "true"
	falseStr      = "false"
)

func InitRedis(redisConfig *config.RedisConfig, logger *zap.Logger) (*RedisManager, error) {
	cfg := redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.Db,
	}

	client := redis.NewClient(&cfg)

	pong, err := client.Ping().Result()
	if err != nil {
		logger.Fatal(
			constants.ERROR_REDIS_CONNECTION_MSG,
			zap.Int32("errorCode", constants.ERROR_REDIS_CONNECTION),
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info(
		constants.INFO_REDIS_CONNECT_SUCCESS,
		zap.String("pong", pong),
	)

	redisManager := RedisManager{
		client: client,
		config: redisConfig,
		logger: logger,
	}

	return &redisManager, err
}

func (rm *RedisManager) Set(key string, bytes []byte, exp time.Duration) error {
	successStr := trueStr
	// time redis op
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.RedisOpDuration.WithLabelValues(rm.config.ServiceLabel, redisOpGetStr, successStr).Observe(v)
	}))
	defer func() {
		timer.ObserveDuration()
	}()

	// call the redis client
	err := rm.client.Set(key, bytes, exp).Err()
	if err != nil {
		rm.logger.Error(
			constants.ERROR_REDIS_SET_MSG,
			zap.String("key", key),
			zap.ByteString("bytes", bytes),
			zap.Error(err),
		)
		// if error, set success to false
		successStr = falseStr
		return errors.Error{constants.ERROR_REDIS_SET, constants.ERROR_REDIS_SET_MSG, err}
	}

	rm.logger.Info(
		constants.INFO_REDIS_SET,
		zap.String("key", key),
		zap.ByteString("bytes", bytes),
		zap.Duration("exp", exp),
	)
	return nil
}

func (rm *RedisManager) Get(key string) ([]byte, error) {
	successStr := trueStr
	// time redis op
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.RedisOpDuration.WithLabelValues(rm.config.ServiceLabel, redisOpGetStr, successStr).Observe(v)
	}))
	defer func() {
		timer.ObserveDuration()
	}()

	// call the redis client
	bytes, err := rm.client.Get(key).Bytes()
	if err != nil {
		if err != redis.Nil {
			// unexpected error occured when getting item
			rm.logger.Error(
				constants.ERROR_REDIS_GET_MSG,
				zap.String("key", key),
				zap.Error(err),
			)
			// set success to false only if unexpected error occured
			successStr = falseStr
			return nil, errors.Error{constants.ERROR_REDIS_GET, constants.ERROR_REDIS_GET_MSG, err}
		}
		rm.logger.Info(
			constants.INFO_REDIS_NOT_FOUND,
			zap.String("key", key),
		)
		// item is not in redis
		return nil, nil
	}

	rm.logger.Info(
		constants.INFO_REDIS_GET,
		zap.String("key", key),
		zap.ByteString("bytes", bytes),
	)
	return bytes, err
}
