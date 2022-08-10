package db

import (
	"fmt"
	"itemService/config"
	constants "itemService/constants"
	errors "itemService/errors"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type RedisManager struct {
	client *redis.Client
	config *config.RedisConfig
	logger *zap.Logger
}

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
	// rm.logger.Info(
	// 	constants.INFO_REDIS_SET,
	// 	zap.String("key", key),
	// 	zap.Any("val", val),
	// 	zap.Duration("exp", exp),
	// )
	// bytes, err := util.MarshalProto(&val)
	// if err != nil {
	// 	rm.logger.Error(
	// 		constants.ERROR_UNMARSHAL_MSG,
	// 		zap.String("key", key),
	// 		zap.Any("val", val),
	// 		zap.Error(err),
	// 	)
	// 	return errors.Error{constants.ERROR_UNMARSHAL, constants.ERROR_UNMARSHAL_MSG, err}
	// }

	err := rm.client.Set(key, bytes, exp).Err()
	if err != nil {
		rm.logger.Error(
			constants.ERROR_REDIS_SET_MSG,
			zap.String("key", key),
			zap.ByteString("bytes", bytes),
			zap.Error(err),
		)
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
	bytes, err := rm.client.Get(key).Bytes()
	if err != nil {
		if err != redis.Nil {
			// unexpected error occured when getting item
			rm.logger.Error(
				constants.ERROR_REDIS_GET_MSG,
				zap.String("key", key),
				zap.Error(err),
			)
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
