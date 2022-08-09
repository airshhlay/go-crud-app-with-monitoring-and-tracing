package server

import (
	"context"
	"fmt"
	"net"
	"userService/config"
	constants "userService/constants"
	"userService/db"
	customErr "userService/errors"
	pb "userService/proto"

	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	handler Handler
	logger  *zap.Logger
}

func (s *Server) StartServer(config *config.Config, dbManager *db.DbManager, logger *zap.Logger) {
	s.handler = Handler{
		config:    config,
		dbManager: dbManager,
		logger:    logger,
	}
	s.logger = logger

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Port))
	if err != nil {
		logger.Fatal(
			constants.ERROR_SERVER_START_FAIL_MSG,
			zap.Int32("errorCode", constants.ERROR_SERVER_START_FAIL),
			zap.Error(err),
		)
	}

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, s)
	logger.Info(
		constants.INFO_SERVER_START_MSG,
		zap.Any("port", listener.Addr()),
	)

	reflection.Register(server)
	err = server.Serve(listener)
	if err != nil {
		logger.Fatal(
			constants.ERROR_SERVER_START_FAIL_MSG,
			zap.Int32("errorCode", constants.ERROR_SERVER_START_FAIL),
			zap.Error(err),
		)
	}
}

func (s *Server) Signup(ctx context.Context, req *pb.SignupReq) (*pb.SignupRes, error) {
	// s.logger.Info(
	// 	constants.INFO_USER_SIGNUP_REQ,
	// 	zap.Any("request", req),
	// )
	// check if a user with the given username exists

	// user does not exist, insert into database
	_, err := s.handler.CreateNewUser(req.Username, req.Password)
	if err != nil {
		v, ok := err.(*customErr.Error)
		if !ok {
			s.logger.Error(
				constants.ERROR_TYPECAST_MSG,
				zap.Error(err),
			)
			return &pb.SignupRes{
				ErrorCode: constants.ERROR_TYPECAST,
			}, nil
		}
		return &pb.SignupRes{
			ErrorCode: v.ErrorCode,
		}, nil
	}

	return &pb.SignupRes{
		ErrorCode: -1,
	}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	var userId int64

	// check if a user with the given username exists
	userId, err := s.handler.VerifyLogin(req.Username, req.Password)
	if err != nil {
		v, ok := err.(*customErr.Error)
		if !ok {
			s.logger.Error(
				constants.ERROR_TYPECAST_MSG,
				zap.Error(err),
			)
			return &pb.LoginRes{
				ErrorCode: constants.ERROR_TYPECAST,
			}, nil
		}
		return &pb.LoginRes{
			ErrorCode: v.ErrorCode,
		}, nil
	}

	return &pb.LoginRes{
		ErrorCode: -1,
		UserId:    userId,
	}, nil
}
