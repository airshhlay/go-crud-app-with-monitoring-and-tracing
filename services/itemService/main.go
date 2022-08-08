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
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	// read in config
	config, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal(
			"Failed to load config",
			zap.Error(err),
		)
	}

	logger.Info("Loaded config")

	// connect to database, get database manager
	dbManager, _ := db.InitDatabase(&config.DbConfig, logger)

	// connect to redis, get redis manager
	redisManager, _ := db.InitRedis(&config.RedisConfig, logger)
	// create server struct
	server := server.Server{}

	// start grpc server
	server.StartServer(config, logger, dbManager, redisManager)
}
