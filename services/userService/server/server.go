package server

import (
	"context"
	"fmt"
	"net"
	"userService/config"
	constants "userService/constants"
	"userService/db"
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
	_, error := s.handler.CreateNewUser(req.Username, req.Password)
	return &pb.SignupRes{
		ErrorCode: error.errorCode,
		ErrorMsg:  error.errorMsg,
	}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	var errorCode int32
	var errorMsg string
	var userId int64

	// check if a user with the given username exists
	userId, error := s.handler.VerifyLogin(req.Username, req.Password)
	if error.errorCode != 0 {
		errorCode = error.errorCode
		errorMsg = error.errorMsg
	}

	return &pb.LoginRes{
		ErrorCode: errorCode,
		ErrorMsg:  errorMsg,
		UserId:    userId,
	}, nil
}
