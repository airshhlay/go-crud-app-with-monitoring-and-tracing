package server

import (
	"context"
	"fmt"
	"net"
	"userService/config"
	"userService/db"
	pb "userService/proto"

	"go.uber.org/zap"

	"log"

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

	listener, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Printf("Error occured when starting the server: %v\n", err)
	}

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, s)
	fmt.Printf("Server listening on port %v\n", listener.Addr())

	reflection.Register(server)
	err = server.Serve(listener)
	if err != nil {
		s.logger.Fatal(
			"Server failed to serve",
			zap.Error(err),
		)
	}
}

func (s *Server) Signup(ctx context.Context, req *pb.SignupReq) (*pb.SignupRes, error) {
	// check if a user with the given username exists
	exists, error := s.handler.CheckUserExists(req.Username)
	if error.errorCode != 0 {
		return &pb.SignupRes{
			ErrorCode: error.errorCode,
			ErrorMsg:  error.errorMsg,
		}, nil
	}

	// if the user does not exist, insert user into database
	if exists {
		_, error = s.handler.CreateNewUser(req.Username, req.Password)
		if error.errorCode != 0 {
			return &pb.SignupRes{
				ErrorCode: error.errorCode,
				ErrorMsg:  error.errorMsg,
			}, nil
		}
	}
	return &pb.SignupRes{}, nil
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
