package main

import (
	"log"
	"userService/config"
	"userService/db"

	"go.uber.org/zap"

	server "userService/server"
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

	// create server struct
	server := server.Server{}

	// start grpc server
	server.StartServer(config, dbManager, logger)
}
