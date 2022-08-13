package client

import (
	"context"
	"fmt"
	config "gateway/config"
	proto "gateway/proto"

	opentracing "github.com/opentracing/opentracing-go"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// ItemServiceClient serves as a wrapper for the grpc client to the item service grpc server
type ItemServiceClient struct {
	logger *zap.Logger
	config *config.GrpcServiceConfig
	client proto.ItemServiceClient
}

// GetItemServiceClient returns the ItemServiceClient
func GetItemServiceClient(logger *zap.Logger, config *config.GrpcServiceConfig) *ItemServiceClient {
	return &ItemServiceClient{
		logger: logger,
		config: config,
	}
}

// StartClient connects to the item service grpc server
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

// AddFav calls the item service's method with the defined AddFavReq
func (i *ItemServiceClient) AddFav(ctx context.Context, req *proto.AddFavReq) (*proto.AddFavRes, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Create")
	defer span.Finish()
	return i.client.AddFav(ctx, req)
}

// DeleteFav calls the item service's method with the defined DeleteFavReq
func (i *ItemServiceClient) DeleteFav(ctx context.Context, req *proto.DeleteFavReq) (*proto.DeleteFavRes, error) {
	return i.client.DeleteFav(ctx, req)
}

// GetFavList calls the item service's method with the defined GetFavList
func (i *ItemServiceClient) GetFavList(ctx context.Context, req *proto.GetFavListReq) (*proto.GetFavListRes, error) {
	return i.client.GetFavList(ctx, req)
}
