package main

import (
	"fmt"
	client "gateway/client"
	constants "gateway/constants"
	controllers "gateway/controllers"
	metrics "gateway/metrics"
	middleware "gateway/middleware"
	routes "gateway/routes"
	jaegerTracer "gateway/tracing"

	opentracing "github.com/opentracing/opentracing-go"

	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"

	config "gateway/config"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// GrpcClients struct holds references to the user service grpc client and item service grpc client
type GrpcClients struct {
	UserServiceClient *client.UserServiceClient
	ItemServiceClient *client.ItemServiceClient
}

func main() {
	logger, err := newLogger()
	if err != nil {
		panic(err)
	}

	config, err := config.LoadConfig(logger)
	if err != nil {
		panic(err)
	}

	// init jaeger
	tracer, closer, err := jaegerTracer.InitJaeger(&config.JaegerConfig, logger)
	if err != nil {
		panic(err)
	}
	opentracing.SetGlobalTracer(tracer)
	logger.Info(constants.InfoJaegerInit)
	defer closer.Close()

	// initialise metrics metrics
	metrics.Init()
	// start the grpc server
	clients := StartGrpcClients(logger, config)
	// start http server
	StartHTTPServer(logger, config, clients)
}

// StartHTTPServer initialise necessary middleware, item service, user service and metrics routes, and starts the HTTP server.
func StartHTTPServer(logger *zap.Logger, config *config.Config, clients *GrpcClients) {
	if config.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	server := gin.New()
	server.SetTrustedProxies([]string{"*"})

	// ignore metrics endpoint when logging
	server.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{config.PrometheusConfig.Endpoint}}))

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	server.Use(gin.Recovery())
	server.Use(middleware.CORSMiddleware(config))

	// prometheus metrics endpoint
	server.GET(config.PrometheusConfig.Endpoint, metrics.PrometheusHandler())

	// Routes for User Service
	userServiceGroup := server.Group(config.HTTPConfig.UserService.URLGroup)
	userServiceController := controllers.NewUserServiceController(&config.HTTPConfig.UserService, logger, clients.UserServiceClient)
	userServiceGroup.Use(middleware.PrometheusMiddleware(config)) // use prometheus middleware
	routes.UserServiceRoutes(userServiceGroup, userServiceController, &config.HTTPConfig.UserService.APIs)

	// Routes for Item Service
	itemServiceGroup := server.Group(config.HTTPConfig.ItemService.URLGroup)
	itemServiceController := controllers.NewItemServiceController(&config.HTTPConfig.ItemService, logger, clients.ItemServiceClient)
	itemServiceGroup.Use(middleware.Authenticate(config.HTTPConfig.UserService.Secret, logger)) // authenticate requests to item service
	itemServiceGroup.Use(middleware.PrometheusMiddleware(config))                               // use prometheus middleware
	routes.ItemServiceRoutes(itemServiceGroup, itemServiceController, &config.HTTPConfig.ItemService.APIs)

	err := server.Run(fmt.Sprintf(":%s", config.Port))
	if err != nil {
		logger.Fatal(
			constants.ErrorServerStartFailMsg,
			zap.Error(err),
		)
		panic(err)
	}

	logger.Info(
		constants.InfoInfoHTTPServerStart,
		zap.String("port", config.Port),
	)
}

// StartGrpcClients starts the grpc client connections to the microservice grpc servers. It returns a reference to the GrpcClients struct.
func StartGrpcClients(logger *zap.Logger, config *config.Config) *GrpcClients {
	generalOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	// user service client
	userServiceClient := client.GetUserServiceClient(logger, &config.GrpcConfig.UserService)
	err := userServiceClient.StartClient(generalOpts)
	if err != nil {
		logger.Fatal(
			constants.ErrorGrpcClientStartFailMsg,
			zap.Error(err),
		)
		panic(err)
	}

	// item service client
	itemServiceClient := client.GetItemServiceClient(logger, &config.GrpcConfig.ItemService)
	err = itemServiceClient.StartClient(generalOpts)
	if err != nil {
		logger.Fatal(
			constants.ErrorGrpcClientStartFailMsg,
			zap.Error(err),
		)
		panic(err)
	}

	return &GrpcClients{
		UserServiceClient: userServiceClient,
		ItemServiceClient: itemServiceClient,
	}
}

func newLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"./log/service.log", "stderr",
	}
	return cfg.Build()
}
