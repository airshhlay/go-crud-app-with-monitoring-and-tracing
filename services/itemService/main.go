package main

import (
	"itemService/config"
	"itemService/db"
	"log"

	"go.uber.org/zap"

	server "itemService/server"
)

func main() {
	// set logger
	logger, err := newLogger()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	// read in config
	config, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal(
			"Failed to load config",
			zap.Error(err),
		)
		panic(err)
	}

	logger.Info("Loaded config")

	// connect to database, get database manager
	dbManager, err := db.InitDatabase(&config.DbConfig, logger)
	if err != nil {
		panic(err)
	}

	// connect to redis, get redis manager
	redisManager, err := db.InitRedis(&config.RedisConfig, logger)
	if err != nil {
		panic(err)
	}
	// create server struct
	server := server.Server{}

	// start grpc server
	server.StartServer(config, logger, dbManager, redisManager)
}

func newLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"./log/service.log", "stderr",
	}
	return cfg.Build()
}
