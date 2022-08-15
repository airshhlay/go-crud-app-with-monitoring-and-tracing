package client

import (
	"context"
	"fmt"
	config "gateway/config"
	proto "gateway/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// UserServiceClient serves as a wrapper for the grpc client to the item service grpc server
type UserServiceClient struct {
	logger *zap.Logger
	config *config.GrpcServiceConfig
	client proto.UserServiceClient
}

// GetUserServiceClient returns the UserServiceClient
func GetUserServiceClient(logger *zap.Logger, config *config.GrpcServiceConfig) *UserServiceClient {
	return &UserServiceClient{
		logger: logger,
		config: config,
	}
}

// StartClient connects to the user service grpc server
func (u *UserServiceClient) StartClient(opts []grpc.DialOption) error {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", u.config.Host, u.config.Port), opts...)
	if err != nil {
		u.logger.Fatal(
			"error_create_grpc_channel",
			zap.String("label", u.config.Label),
			zap.String("host", u.config.Host),
			zap.String("port", u.config.Port),
			zap.Error(err),
		)
		panic(err)
	}

	// defer conn.Close()

	client := proto.NewUserServiceClient(conn)
	u.client = client

	u.logger.Info(
		"info_grpc_client_start",
		zap.String("label", u.config.Label),
		zap.String("host", u.config.Host),
		zap.String("port", u.config.Port),
	)

	return err
}

// Signup calls the user service's method with the defined SignupReq
func (u *UserServiceClient) Signup(ctx context.Context, req *proto.SignupReq) (*proto.SignupRes, error) {
	// get signup request
	return u.client.Signup(ctx, req)
}

// Login calls the user service's method with the defined LoginReq
func (u *UserServiceClient) Login(ctx context.Context, req *proto.LoginReq) (*proto.LoginRes, error) {
	res, err := u.client.Login(ctx, req)
	if err != nil {
		u.logger.Error(
			"error_userservice_login",
			zap.Error(err),
		)
	}
	return res, err
}
