package main

import (
	"fmt"
	client "gateway/client"
	controllers "gateway/controllers"
	"gateway/middleware"
	routes "gateway/routes"

	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"

	config "gateway/config"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type GrpcClients struct {
	UserServiceClient *client.UserServiceClient
	ItemServiceClient *client.ItemServiceClient
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	config, err := config.LoadConfig(logger)
	if err != nil {
		panic(err)
	}

	// start the grpc server
	clients := StartGrpcClients(logger, config)
	// start http server
	StartHttpServer(logger, config, clients)
}

func StartHttpServer(logger *zap.Logger, config *config.Config, clients *GrpcClients) {
	if config.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	server := gin.Default()
	server.SetTrustedProxies([]string{"*"})
	server.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	server.Use(gin.Recovery())
	// server.Use(cors.Default())
	server.Use(middleware.CORSMiddleware())

	// health check
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Routes for User Service
	userServiceGroup := server.Group(config.HttpConfig.UserService.UrlGroup)
	userServiceController := controllers.NewUserServiceController(&config.HttpConfig.UserService, logger, clients.UserServiceClient)
	routes.UserServiceRoutes(userServiceGroup, userServiceController, &config.HttpConfig.UserService.Apis)

	// Routes for Item Service
	itemServiceGroup := server.Group(config.HttpConfig.ItemService.UrlGroup)
	itemServiceController := controllers.NewItemServiceController(&config.HttpConfig.ItemService, logger, clients.ItemServiceClient)
	itemServiceGroup.Use(middleware.Authenticate(config.HttpConfig.UserService.Secret))
	routes.ItemServiceRoutes(itemServiceGroup, itemServiceController, &config.HttpConfig.ItemService.Apis)

	// server.Use(cors.Default())
	// server.Use(middleware.CORSMiddleware())
	// server.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:3000"},
	// 	AllowMethods:     []string{"GET, POST, PUT, DELETE"},
	// 	AllowHeaders:     []string{"Origin"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))
	err := server.Run(fmt.Sprintf(":%s", config.Port))
	if err != nil {
		logger.Fatal(
			"error_server_failure",
			zap.Error(err),
		)
		panic(err)
	}

	logger.Info(
		"info_http_rest_server_start",
		zap.String("port", config.Port),
	)
}

func StartGrpcClients(logger *zap.Logger, config *config.Config) *GrpcClients {
	generalOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	// user service client
	userServiceClient := client.GetUserServiceClient(logger, &config.GrpcConfig.UserService)
	// var userServiceOpts []grpc.DialOption
	userServiceOpts := generalOpts
	err := userServiceClient.StartClient(userServiceOpts)
	if err != nil {
		logger.Fatal(
			"error_server_failure",
			zap.Error(err),
		)
		panic(err)
	}

	// var itemServiceOpts []grpc.DialOption
	itemServiceOpts := generalOpts
	itemServiceClient := client.GetItemServiceClient(logger, &config.GrpcConfig.ItemService)
	err = itemServiceClient.StartClient(itemServiceOpts)
	if err != nil {
		logger.Fatal(
			"error_server_failure",
			zap.Error(err),
		)
		panic(err)
	}

	return &GrpcClients{
		UserServiceClient: userServiceClient,
		ItemServiceClient: itemServiceClient,
	}
}
