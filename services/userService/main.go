package main

import (
	"context"
	"log"
	"net"
	"userService/config"
	"userService/db"
	pb "userService/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) Signup(ctx context.Context, in *pb.SignupReq) (*pb.SignupRes, error) {
	return &pb.SignupRes{
		ErrorCode: 1,
		ErrorMsg:  "No error",
	}, nil
}

func main() {
	// read in config
	config, _ := config.LoadConfig()
	log.Printf("Read config %v", config)

	// connect to database, get database manager
	dbManager, _ := db.Init(&config.DbConfig)
	dbManager.DoSomething()

	listener, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Printf("Error occured when starting the server: %v\n", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})
	log.Printf("Server listening at %v\n", listener.Addr())

	reflection.Register(s)
	err = s.Serve(listener)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
