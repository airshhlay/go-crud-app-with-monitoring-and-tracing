package db

import (
	"context"
	"fmt"
	"itemService/config"
	constants "itemService/constants"
	errors "itemService/errors"
	metrics "itemService/metrics"
	"itemService/tracing"
	"time"

	ot "github.com/opentracing/opentracing-go"

	"github.com/prometheus/client_golang/prometheus"

	redis "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const (
	redisSet = "redis.Set"
	redisGet = "redis.Get"
)

// RedisManager is a struct containing a reference to the redis client, logger, and the redis config
type RedisManager struct {
	client *redis.Client
	config *config.RedisConfig
	logger *zap.Logger
}

// InitRedis creates the redis client, initialises and tests the connection
func InitRedis(redisConfig *config.RedisConfig, logger *zap.Logger) (*RedisManager, error) {
	ctx := context.Background()
	cfg := redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.Db,
	}

	client := redis.NewClient(&cfg)

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		logger.Fatal(
			constants.ErrorRedisConnectionMsg,
			zap.Int32(constants.ErrorCode, constants.ErrorRedisConnection),
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info(
		constants.InfoRedisConnectSuccess,
		zap.String("pong", pong),
	)

	redisManager := RedisManager{
		client: client,
		config: redisConfig,
		logger: logger,
	}

	return &redisManager, err
}

// Set takes a key of type string and a byte array as a value. exp is used to define an expiry.
func (rm *RedisManager) Set(ctx context.Context, key string, bytes []byte, exp time.Duration) error {
	// start tracing span from context
	span, ctx := ot.StartSpanFromContext(ctx, redisSet)
	statement := fmt.Sprintf(tracing.DatabaseStatementRedisSet, key, bytes)
	rm.addSpanTags(span, statement)
	defer span.Finish()
	successStr := constants.True
	// time redis op
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.RedisOpDuration.WithLabelValues(rm.config.ServiceLabel, constants.Get, successStr).Observe(v)
	}))
	defer func() {
		timer.ObserveDuration()
	}()

	// call the redis client
	err := rm.client.Set(ctx, key, bytes, exp).Err()
	if err != nil {
		rm.logger.Error(
			constants.ErrorRedisSetMsg,
			zap.String(constants.Key, key),
			zap.ByteString(constants.Bytes, bytes),
			zap.Error(err),
		)
		// if error, set success to false
		successStr = constants.False
		return errors.Error{constants.ErrorRedisSet, constants.ErrorRedisSetMsg, err}
	}

	rm.logger.Info(
		constants.InfoRedisSet,
		zap.String(constants.Key, key),
		zap.ByteString(constants.Bytes, bytes),
		zap.Duration(constants.Exp, exp),
	)
	return nil
}

// Get takes a key and returns its associated value in bytes.
func (rm *RedisManager) Get(ctx context.Context, key string) ([]byte, error) {
	// start tracing span from context
	span, ctx := ot.StartSpanFromContext(ctx, redisGet)
	statement := fmt.Sprintf(tracing.DatabaseStatementRedisGet, key)
	rm.addSpanTags(span, statement)
	defer span.Finish()
	successStr := constants.True
	// time redis op
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.RedisOpDuration.WithLabelValues(rm.config.ServiceLabel, constants.Get, successStr).Observe(v)
	}))
	defer func() {
		timer.ObserveDuration()
	}()

	// call the redis client
	bytes, err := rm.client.Get(ctx, key).Bytes()
	if err != nil {
		if err != redis.Nil {
			// unexpected error occured when getting item
			rm.logger.Error(
				constants.ErrorRedisGetMsg,
				zap.String(constants.Key, key),
				zap.Error(err),
			)
			// set success to false only if unexpected error occured
			successStr = constants.False
			return nil, errors.Error{constants.ErrorRedisGet, constants.ErrorRedisGetMsg, err}
		}
		rm.logger.Info(
			constants.InfoItemNotInRedis,
			zap.String(constants.Key, key),
		)
		// item is not in redis
		return nil, nil
	}

	rm.logger.Info(
		constants.InfoRedisGet,
		zap.String(constants.Key, key),
		zap.ByteString(constants.Bytes, bytes),
	)
	return bytes, err
}

func (rm *RedisManager) addSpanTags(span ot.Span, statement string) {
	span.SetTag(tracing.DatabaseType, tracing.DatabaseTypeRedis)
	span.SetTag(tracing.DatabaseInstance, rm.config.Db)
	span.SetTag(tracing.DatabaseUser, rm.config.Host)
	span.SetTag(tracing.DatabaseStatement, statement)
	span.SetTag(tracing.Component, tracing.ComponentDB)
}
