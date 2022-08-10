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
	logger, err := NewLogger()
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
	dbManager, err := db.InitDatabase(&config.DbConfig, logger)
	if err != nil {
		panic(err)
	}

	// create server struct
	server := server.Server{}

	// start grpc server
	server.StartServer(config, dbManager, logger)
}

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"./log/service.log", "stderr",
	}
	return cfg.Build()
}
