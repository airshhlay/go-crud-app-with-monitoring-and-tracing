package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"userService/config"
	constants "userService/constants"
	"userService/db"

	customErr "userService/errors"
	metrics "userService/metrics"
	pb "userService/proto"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	SIGNUP = "signup"
	LOGIN  = "login"
)

// var (
// 	// Create a metrics registry.
// 	reg = prometheus.NewRegistry()
// 	// Create some standard server metrics.
// 	grpcMetrics = grpc_prometheus.NewServerMetrics()
// 	// Create a customized counter metric.
// 	customizedCounterMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
// 		Name: "demo_server_say_hello_method_handle_count",
// 		Help: "Total number of RPCs handled on the server.",
// 	}, []string{"name"})
// )

// func init() {
// 	// Register standard server metrics and customized metrics to registry.
// 	reg.MustRegister(grpcMetrics, customizedCounterMetric)
// 	customizedCounterMetric.WithLabelValues("Test")
// }

type Server struct {
	pb.UnimplementedUserServiceServer
	handler Handler
	logger  *zap.Logger
	config  *config.Config
}

func (s *Server) StartServer(config *config.Config, dbManager *db.DbManager, logger *zap.Logger) {
	s.handler = Handler{
		config:    config,
		dbManager: dbManager,
		logger:    logger,
	}
	s.logger = logger
	s.config = config

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Port))
	if err != nil {
		logger.Fatal(
			constants.ERROR_SERVER_START_FAIL_MSG,
			zap.Int32("errorCode", constants.ERROR_SERVER_START_FAIL),
			zap.Error(err),
		)
	}
	defer listener.Close()

	// Create a HTTP server for prometheus.
	http.Handle(config.PrometheusConfig.Endpoint, promhttp.HandlerFor(metrics.Reg, promhttp.HandlerOpts{}))
	// http.Handle(config.PrometheusConfig.Endpoint, promhttp.Handler())
	// httpServer := &http.Server{Handler: , Addr: fmt.Sprintf("%s:%s", config.PrometheusConfig.Host, config.PrometheusConfig.Port)}

	// Create a gRPC Server with gRPC interceptor.
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(metrics.GrpcMetrics.StreamServerInterceptor()),
		grpc.UnaryInterceptor(metrics.GrpcMetrics.UnaryServerInterceptor()),
	)
	if err != nil {
		logger.Error(
			constants.ERROR_PROM_INIT_CUSTOM_METRICS_MSG,
			zap.Error(err),
		)
	}

	// register service
	pb.RegisterUserServiceServer(grpcServer, s)

	// initialize all metrics.
	metrics.GrpcMetrics.InitializeMetrics(grpcServer)

	// start http server for prometheus
	go func() {
		logger.Info("Starting http server", zap.String("host", config.PrometheusConfig.Host), zap.String("port", config.PrometheusConfig.Port))
		// err := httpServer.ListenAndServe()
		http.ListenAndServe(fmt.Sprintf(":%s", config.PrometheusConfig.Port), nil)
		if err != nil {
			logger.Fatal(constants.ERROR_PROM_HTTP_SERVER_MSG,
				zap.Error(err))
		} else {
			logger.Info(constants.INFO_PROM_SERVER_START_MSG)
		}
	}()

	// start grpc server
	reflection.Register(grpcServer)
	err = grpcServer.Serve(listener)
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

func (s *Server) Signup(ctx context.Context, req *pb.SignupReq) (*pb.SignupRes, error) {
	errorCodeStr := "-1"
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.DatabaseOpDuration.WithLabelValues(s.config.ServiceLabel, SIGNUP, GET_USER_BY_USERNAME_OP, errorCodeStr).Observe(v)
	}))

	// observe duration at the end of this function
	defer func() {
		timer.ObserveDuration()
	}()
	// user does not exist, insert into database
	_, err := s.handler.CreateNewUser(req.Username, req.Password)
	if err != nil {
		v, ok := err.(*customErr.Error)
		if !ok {
			s.logger.Error(constants.ERROR_TYPECAST_MSG, zap.Error(err))
			errorCodeStr = strconv.Itoa(constants.ERROR_TYPECAST)
			return &pb.SignupRes{
				ErrorCode: constants.ERROR_TYPECAST,
			}, nil
		}
		errorCodeStr = strconv.Itoa(int(v.ErrorCode))
		return &pb.SignupRes{
			ErrorCode: v.ErrorCode,
		}, nil
	}

	return &pb.SignupRes{
		ErrorCode: -1,
	}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	errorCodeStr := "-1"
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.DatabaseOpDuration.WithLabelValues(s.config.ServiceLabel, LOGIN, GET_USER_BY_USERNAME_OP, errorCodeStr).Observe(v)
	}))

	// observe duration at the end of this function
	defer func() {
		timer.ObserveDuration()
	}()

	var userId int64

	// check if a user with the given username exists
	userId, err := s.handler.VerifyLogin(req.Username, req.Password)
	if err != nil {
		v, ok := err.(*customErr.Error)
		if !ok {
			s.logger.Error(constants.ERROR_TYPECAST_MSG, zap.Error(err))
			errorCodeStr = strconv.Itoa(constants.ERROR_TYPECAST)
			return &pb.LoginRes{
				ErrorCode: constants.ERROR_TYPECAST,
			}, nil
		}
		errorCodeStr = strconv.Itoa(int(v.ErrorCode))
		return &pb.LoginRes{
			ErrorCode: v.ErrorCode,
		}, nil
	}

	return &pb.LoginRes{
		ErrorCode: -1,
		UserId:    userId,
	}, nil
}
