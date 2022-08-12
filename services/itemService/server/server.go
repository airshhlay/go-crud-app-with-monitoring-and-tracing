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

const (
	deleteFavStr  = "deleteFav"
	addFavStr     = "addFav"
	getFavListStr = "getFavList"
)

// Server struct contains a reference to the handler. Used to start the grpc server.
type Server struct {
	pb.UnimplementedItemServiceServer
	handler Handler
	logger  *zap.Logger
	config  *config.Config
}

// StartServer initialises the prometheus metrics, starts the HTTP server for the prometheus endpoint and starts the GRPC server.
func (s *Server) StartServer(config *config.Config, logger *zap.Logger, dbManager *db.DbManager, redisManager *db.RedisManager) {
	s.handler = Handler{
		config:       config,
		dbManager:    dbManager,
		redisManager: redisManager,
		logger:       logger,
	}
	s.logger = logger
	s.config = config

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Port))
	if err != nil {
		s.logger.Fatal(
			constants.ERROR_SERVER_START_FAIL_MSG,
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
		logger.Info("Starting http server", zap.String("host", config.PrometheusConfig.Host), zap.String("port", config.PrometheusConfig.Port))
		http.ListenAndServe(fmt.Sprintf(":%s", config.PrometheusConfig.Port), nil)
		if err != nil {
			logger.Fatal(constants.ErrorPromHTTPServerMsg,
				zap.Error(err))
		} else {
			logger.Info(constants.INFO_PROM_SERVER_START_MSG)
		}
	}()

	reflection.Register(grpcServer)
	err = grpcServer.Serve(listener)
	if err != nil {
		logger.Fatal(
			constants.ERROR_SERVER_START_FAIL_MSG,
			zap.Error(err),
		)
	}
	if err != nil {
		logger.Fatal(
			constants.ERROR_SERVER_START_FAIL_MSG,
			zap.Int32("errorCode", constants.ERROR_SERVER_START_FAIL),
			zap.Error(err),
		)
	} else {
		logger.Info(
			constants.INFO_SERVER_START_MSG,
			zap.Any("port", listener.Addr()),
		)
	}
}

// DeleteFav is the implementation of the grpc server service, as defined in service.proto
func (s *Server) DeleteFav(ctx context.Context, req *pb.DeleteFavReq) (*pb.DeleteFavRes, error) {
	errorCodeStr := "-1"
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.RequestDuration.WithLabelValues(s.config.ServiceLabel, deleteFavStr, errorCodeStr).Observe(v)
	}))
	// observe duration at the end of this function
	defer func() {
		timer.ObserveDuration()
	}()

	err := s.handler.DeleteFavourite(req.UserId, req.ItemId, req.ShopId)
	if err != nil {
		v, ok := err.(*customErr.Error)
		if !ok {
			s.logger.Error(constants.ERROR_TYPECAST_MSG, zap.Error(err))
			errorCodeStr = strconv.Itoa(constants.ERROR_TYPECAST)
			return &pb.DeleteFavRes{
				ErrorCode: constants.ERROR_TYPECAST,
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

func (s *Server) AddFav(ctx context.Context, req *pb.AddFavReq) (*pb.AddFavRes, error) {
	errorCodeStr := "-1"
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.RequestDuration.WithLabelValues(s.config.ServiceLabel, addFavStr, errorCodeStr).Observe(v)
	}))
	// observe duration at the end of this function
	defer func() {
		timer.ObserveDuration()
	}()

	item, err := s.handler.AddItemToUserFavList(req.ItemId, req.ShopId, req.UserId)

	if err != nil {
		v, ok := err.(*customErr.Error)
		if !ok {
			s.logger.Error(
				constants.ERROR_TYPECAST_MSG,
				zap.Error(err),
			)
			errorCodeStr = strconv.Itoa(constants.ERROR_TYPECAST)
			return &pb.AddFavRes{
				ErrorCode: constants.ERROR_TYPECAST,
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

func (s *Server) GetFavList(ctx context.Context, req *pb.GetFavListReq) (*pb.GetFavListRes, error) {
	errorCodeStr := "-1"
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.RequestDuration.WithLabelValues(s.config.ServiceLabel, getFavListStr, errorCodeStr).Observe(v)
	}))
	// observe duration at the end of this function
	defer func() {
		timer.ObserveDuration()
	}()

	items, totalPages, err := s.handler.GetUserFavourites(req.UserId, req.Page)
	if err != nil {
		v, ok := err.(*customErr.Error)
		if !ok {
			s.logger.Error(constants.ERROR_TYPECAST_MSG, zap.Error(err))
			errorCodeStr = strconv.Itoa(constants.ERROR_TYPECAST)
			return &pb.GetFavListRes{
				ErrorCode: constants.ERROR_TYPECAST,
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
