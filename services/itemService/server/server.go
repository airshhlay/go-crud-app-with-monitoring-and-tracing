package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	config "itemService/config"
	constants "itemService/constants"
	db "itemService/db"
	customErr "itemService/errors"
	metrics "itemService/metrics"
	pb "itemService/proto"
	"net"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/prometheus/client_golang/prometheus"

	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server struct contains a reference to the handler. Used to start the grpc server.
type Server struct {
	pb.UnimplementedItemServiceServer
	handler Handler
	logger  *zap.Logger
	config  *config.Config
}

// StartServer initialises the prometheus metrics, starts the HTTP server for the prometheus endpoint and starts the GRPC server.
func (s *Server) StartServer(config *config.Config, logger *zap.Logger, dbManager *db.DatabaseManager, redisManager *db.RedisManager) {
	s.handler = Handler{
		config:       config,
		dbManager:    dbManager,
		redisManager: redisManager,
		logger:       logger,
	}
	s.logger = logger
	s.config = config

	listener, err := net.Listen(constants.TCP, fmt.Sprintf(":%s", config.Port))
	if err != nil {
		s.logger.Fatal(
			constants.ErrorServerStartFailMsg,
			zap.Error(err),
		)
	}
	defer listener.Close()

	// Create a HTTP server for prometheus.
	http.Handle(config.PrometheusConfig.Endpoint, promhttp.HandlerFor(metrics.Reg, promhttp.HandlerOpts{}))

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(metrics.GrpcMetrics.StreamServerInterceptor()),
		grpc.UnaryInterceptor(metrics.GrpcMetrics.UnaryServerInterceptor()),
	)
	if err != nil {
		logger.Error(
			constants.ErrorPromInitCustomMetricsMsg,
			zap.Error(err),
		)
	}

	// register service
	pb.RegisterItemServiceServer(grpcServer, s)

	// initialize all metrics.
	metrics.GrpcMetrics.InitializeMetrics(grpcServer)

	// start http server for prometheus
	go func() {
		logger.Info(constants.InfoHTTPServerStart, zap.String(constants.Host, config.PrometheusConfig.Host), zap.String(constants.Port, config.PrometheusConfig.Port))
		http.ListenAndServe(fmt.Sprintf(":%s", config.PrometheusConfig.Port), nil)
		if err != nil {
			logger.Fatal(constants.ErrorPromHTTPServerMsg,
				zap.Error(err))
		} else {
			logger.Info(constants.InfoPromServerStart)
		}
	}()

	reflection.Register(grpcServer)
	err = grpcServer.Serve(listener)
	if err != nil {
		logger.Fatal(
			constants.ErrorServerStartFailMsg,
			zap.Error(err),
		)
	}
	if err != nil {
		logger.Fatal(
			constants.ErrorServerStartFailMsg,
			zap.Int32(constants.ErrorCode, constants.ErrorServerStartFail),
			zap.Error(err),
		)
	} else {
		logger.Info(
			constants.InfoGRPCServerStart,
			zap.Any(constants.Port, listener.Addr()),
		)
	}
}

// DeleteFav implements the grpc service method, as defined in service.proto
func (s *Server) DeleteFav(ctx context.Context, req *pb.DeleteFavReq) (*pb.DeleteFavRes, error) {
	errorCodeStr := constants.NilErrorCode
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.RequestDuration.WithLabelValues(s.config.ServiceLabel, constants.DeleteFav, errorCodeStr).Observe(v)
	}))
	// observe duration at the end of this function
	defer func() {
		timer.ObserveDuration()
	}()

	err := s.handler.DeleteFavourite(req.UserID, req.ItemID, req.ShopID)
	if err != nil {
		v, ok := err.(*customErr.Error)
		if !ok {
			s.logger.Error(constants.ErrorTypecastMsg, zap.Error(err))
			errorCodeStr = strconv.Itoa(constants.ErrorTypecast)
			return &pb.DeleteFavRes{
				ErrorCode: constants.ErrorTypecast,
			}, nil
		}
		errorCodeStr = strconv.Itoa(int(v.ErrorCode))
		return &pb.DeleteFavRes{
			ErrorCode: v.ErrorCode,
		}, nil
	}
	return &pb.DeleteFavRes{
		ErrorCode: -1,
	}, nil
}

// AddFav implements the grpc service method, as defined in service.proto
func (s *Server) AddFav(ctx context.Context, req *pb.AddFavReq) (*pb.AddFavRes, error) {
	errorCodeStr := constants.NilErrorCode
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.RequestDuration.WithLabelValues(s.config.ServiceLabel, constants.AddFav, errorCodeStr).Observe(v)
	}))
	// observe duration at the end of this function
	defer func() {
		timer.ObserveDuration()
	}()

	item, err := s.handler.AddItemToUserFavList(req.ItemID, req.ShopID, req.UserID)

	if err != nil {
		v, ok := err.(*customErr.Error)
		if !ok {
			s.logger.Error(
				constants.ErrorTypecastMsg,
				zap.Error(err),
			)
			errorCodeStr = strconv.Itoa(constants.ErrorTypecast)
			return &pb.AddFavRes{
				ErrorCode: constants.ErrorTypecast,
			}, nil
		}
		errorCodeStr = strconv.Itoa(int(v.ErrorCode))
		return &pb.AddFavRes{
			ErrorCode: v.ErrorCode,
		}, nil
	}
	return &pb.AddFavRes{
		ErrorCode: -1,
		Item:      item,
	}, nil
}

// GetFavList implements the grpc service method, as defined in service.proto
func (s *Server) GetFavList(ctx context.Context, req *pb.GetFavListReq) (*pb.GetFavListRes, error) {
	errorCodeStr := constants.NilErrorCode
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.RequestDuration.WithLabelValues(s.config.ServiceLabel, constants.GetFavList, errorCodeStr).Observe(v)
	}))
	// observe duration at the end of this function
	defer func() {
		timer.ObserveDuration()
	}()

	items, totalPages, err := s.handler.GetUserFavourites(req.UserID, req.Page)
	if err != nil {
		v, ok := err.(*customErr.Error)
		if !ok {
			s.logger.Error(constants.ErrorTypecastMsg, zap.Error(err))
			errorCodeStr = strconv.Itoa(constants.ErrorTypecast)
			return &pb.GetFavListRes{
				ErrorCode: constants.ErrorTypecast,
			}, nil
		}
		errorCodeStr = strconv.Itoa(int(v.ErrorCode))
		return &pb.GetFavListRes{
			ErrorCode: v.ErrorCode,
		}, nil
	}

	return &pb.GetFavListRes{
		ErrorCode:  -1,
		Items:      items,
		TotalPages: totalPages,
	}, nil
}
