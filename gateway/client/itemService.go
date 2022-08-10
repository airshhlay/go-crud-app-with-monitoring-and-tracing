package client

import (
	"context"
	"fmt"
	config "gateway/config"
	proto "gateway/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ItemServiceClient struct {
	logger *zap.Logger
	config *config.GrpcServiceConfig
	client proto.ItemServiceClient
}

func GetItemServiceClient(logger *zap.Logger, config *config.GrpcServiceConfig) *ItemServiceClient {
	return &ItemServiceClient{
		logger: logger,
		config: config,
	}
}

func (i *ItemServiceClient) StartClient(opts []grpc.DialOption) error {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", i.config.Host, i.config.Port), opts...)
	if err != nil {
		i.logger.Fatal(
			"error_create_grpc_channel",
			zap.String("label", i.config.Label),
			zap.String("host", i.config.Host),
			zap.String("port", i.config.Port),
			zap.Error(err),
		)
		panic(err)
	}

	// defer conn.Close()

	client := proto.NewItemServiceClient(conn)
	i.client = client

	i.logger.Info(
		"info_grpc_client_start",
		zap.String("label", i.config.Label),
		zap.String("host", i.config.Host),
		zap.String("port", i.config.Port),
	)

	return err
}

func (i *ItemServiceClient) AddFav(ctx context.Context, req *proto.AddFavReq) (*proto.AddFavRes, error) {
	return i.client.AddFav(ctx, req)
}

func (i *ItemServiceClient) DeleteFav(ctx context.Context, req *proto.DeleteFavReq) (*proto.DeleteFavRes, error) {
	return i.client.DeleteFav(ctx, req)
}

func (i *ItemServiceClient) GetFavList(ctx context.Context, req *proto.GetFavListReq) (*proto.GetFavListRes, error) {
	return i.client.GetFavList(ctx, req)
}
