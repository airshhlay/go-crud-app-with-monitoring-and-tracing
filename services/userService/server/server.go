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
	"userService/tracing"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	otgrpc "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"

	customErr "userService/errors"
	metrics "userService/metrics"
	pb "userService/proto"

	ot "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	grpcSignup = "userservice.server.Signup"
	grpcLogin  = "userservice.server.Login"
)

// Server struct contains a reference to the handler. Used to start the grpc server.
type Server struct {
	pb.UnimplementedUserServiceServer
	handler Handler
	logger  *zap.Logger
	config  *config.Config
}

// StartServer initialises the prometheus metrics, starts the HTTP server for the prometheus endpoint and starts the GRPC server.
func (s *Server) StartServer(config *config.Config, dbManager *db.DatabaseManager, logger *zap.Logger, tracer ot.Tracer) {
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
			constants.ErrorServerStartFailMsg,
			zap.Int32(constants.ErrorCode, constants.ErrorServerStartFail),
			zap.Error(err),
		)
	}
	defer listener.Close()

	// Create a HTTP server for prometheus.
	http.Handle(config.PrometheusConfig.Endpoint, promhttp.HandlerFor(metrics.Reg, promhttp.HandlerOpts{}))

	// Create a gRPC Server with gRPC interceptor.
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(
			metrics.GrpcMetrics.StreamServerInterceptor(),
		),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				metrics.GrpcMetrics.UnaryServerInterceptor(),
				otgrpc.OpenTracingServerInterceptor(tracer),
			),
		),
	)
	if err != nil {
		logger.Error(
			constants.ErrorPromInitCustomMetricsMsg,
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
			logger.Fatal(constants.ErrorPromHTTPServerMsg,
				zap.Error(err))
		} else {
			logger.Info(constants.InfoPromServerStart)
		}
	}()

	// start grpc server
	reflection.Register(grpcServer)
	err = grpcServer.Serve(listener)
	if err != nil {
		logger.Fatal(
			constants.ErrorServerStartFailMsg,
			zap.Int32(constants.ErrorCode, constants.ErrorServerStartFail),
			zap.Error(err),
		)
	} else {
		logger.Info(
			constants.InfoServerStart,
			zap.Any("port", listener.Addr()),
		)
	}
}

// Signup is the implementation of the grpc server service, as defined in service.proto
func (s *Server) Signup(ctx context.Context, req *pb.SignupReq) (*pb.SignupRes, error) {
	// start tracing span from context
	span, ctx := ot.StartSpanFromContext(ctx, grpcSignup)
	s.addSpanTags(span)
	defer span.Finish()

	errorCodeStr := "-1"
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.RequestDuration.WithLabelValues(s.config.ServiceLabel, constants.Signup, errorCodeStr).Observe(v)
	}))

	// observe duration at the end of this function
	defer func() {
		timer.ObserveDuration()
	}()
	// user does not exist, insert into database
	_, err := s.handler.CreateNewUser(ctx, req.Username, req.Password)
	if err != nil {
		v, ok := err.(*customErr.Error)
		if !ok {
			s.logger.Error(constants.ErrorTypecastMsg, zap.Error(err))
			errorCodeStr = strconv.Itoa(constants.ErrorTypecast)
			return &pb.SignupRes{
				ErrorCode: constants.ErrorTypecast,
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

// Login is the implementation of the grpc server service, as defined in service.proto
func (s *Server) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	// start tracing span from context
	span, ctx := ot.StartSpanFromContext(ctx, grpcSignup)
	s.addSpanTags(span)
	defer span.Finish()

	errorCodeStr := "-1"
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.RequestDuration.WithLabelValues(s.config.ServiceLabel, constants.Login, errorCodeStr).Observe(v)
	}))

	// observe duration at the end of this function
	defer func() {
		timer.ObserveDuration()
	}()

	var userID int64

	// check if a user with the given username exists
	userID, err := s.handler.VerifyLogin(ctx, req.Username, req.Password)
	if err != nil {
		v, ok := err.(*customErr.Error)
		if !ok {
			s.logger.Error(constants.ErrorTypecastMsg, zap.Error(err))
			errorCodeStr = strconv.Itoa(constants.ErrorTypecast)
			return &pb.LoginRes{
				ErrorCode: constants.ErrorTypecast,
			}, nil
		}
		errorCodeStr = strconv.Itoa(int(v.ErrorCode))
		return &pb.LoginRes{
			ErrorCode: v.ErrorCode,
		}, nil
	}

	return &pb.LoginRes{
		ErrorCode: -1,
		UserID:    userID,
	}, nil
}

func (s *Server) addSpanTags(span ot.Span) {
	span.SetTag(tracing.SpanKind, tracing.SpanKindServer)
	span.SetTag(tracing.Component, tracing.ComponentServer)
}
