package main

import (
	"log"
	"userService/config"
	"userService/constants"
	"userService/db"

	"go.uber.org/zap"

	server "userService/server"
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
		logger.Fatal(constants.ErrorLoadConfigFailMsg, zap.Error(err))
		panic(err)
	}

	logger.Info(constants.InfoConfigLoaded)

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

func newLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"./log/service.log", "stderr",
	}
	return cfg.Build()
}
