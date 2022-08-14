package client

import (
	"context"
	"fmt"
	config "gateway/config"
	"gateway/constants"
	proto "gateway/proto"
	"gateway/tracing"

	ot "github.com/opentracing/opentracing-go"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	addFavClient     = "gateway.AddFavClient"
	deleteFavClient  = "gateway.DeleteFavClient"
	getFavListClient = "gateway.GetFavListClient"
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
			constants.ErrorCreateGRPCChannelMsg,
			zap.String(constants.Label, i.config.Label),
			zap.String(constants.Host, i.config.Host),
			zap.String(constants.Port, i.config.Port),
			zap.Error(err),
		)
		return err
	}

	// defer conn.Close()

	client := proto.NewItemServiceClient(conn)
	i.client = client

	i.logger.Info(
		constants.InfoGRPCClientStart,
		zap.String(constants.Label, i.config.Label),
		zap.String(constants.Host, i.config.Host),
		zap.String(constants.Port, i.config.Port),
	)

	return err
}

// AddFav calls the item service's method with the defined AddFavReq
func (i *ItemServiceClient) AddFav(ctx context.Context, req *proto.AddFavReq) (*proto.AddFavRes, error) {
	// start span from context
	span, ctx := ot.StartSpanFromContext(ctx, signupClient)
	i.addSpanTags(span)
	defer span.Finish()

	return i.client.AddFav(ctx, req)
}

// DeleteFav calls the item service's method with the defined DeleteFavReq
func (i *ItemServiceClient) DeleteFav(ctx context.Context, req *proto.DeleteFavReq) (*proto.DeleteFavRes, error) {
	// start span from context
	span, ctx := ot.StartSpanFromContext(ctx, signupClient)
	i.addSpanTags(span)
	defer span.Finish()

	return i.client.DeleteFav(ctx, req)
}

// GetFavList calls the item service's method with the defined GetFavList
func (i *ItemServiceClient) GetFavList(ctx context.Context, req *proto.GetFavListReq) (*proto.GetFavListRes, error) {
	// start span from context
	span, ctx := ot.StartSpanFromContext(ctx, signupClient)
	i.addSpanTags(span)
	defer span.Finish()

	return i.client.GetFavList(ctx, req)
}

func (i *ItemServiceClient) addSpanTags(span ot.Span) {
	span.SetTag(tracing.SpanKind, tracing.SpanKindClient)
	span.SetTag(tracing.Component, tracing.ComponentGrpc)
	span.SetTag(tracing.PeerService, tracing.PeerServiceItemService)
}
