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
	component    = "gin"
	peerService  = "gateway"
	spanKind     = "client"
	signupClient = "gateway.SignupClient"
	loginClient  = "gateway.LoginClient"
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
			constants.ErrorCreateGRPCChannelMsg,
			zap.String(constants.Label, u.config.Label),
			zap.String(constants.Host, u.config.Host),
			zap.String(constants.Port, u.config.Port),
			zap.Error(err),
		)
		return err
	}

	// defer conn.Close()

	client := proto.NewUserServiceClient(conn)
	u.client = client

	u.logger.Info(
		constants.InfoGRPCClientStart,
		zap.String(constants.Label, u.config.Label),
		zap.String(constants.Host, u.config.Host),
		zap.String(constants.Port, u.config.Port),
	)

	return err
}

// Signup calls the user service's method with the defined SignupReq
func (u *UserServiceClient) Signup(ctx context.Context, req *proto.SignupReq) (*proto.SignupRes, error) {
	// start span from context
	span, ctx := ot.StartSpanFromContext(ctx, signupClient)
	u.addSpanTags(span)
	defer span.Finish()

	// send the request to user service
	return u.client.Signup(ctx, req)
}

// Login calls the user service's method with the defined LoginReq
func (u *UserServiceClient) Login(ctx context.Context, req *proto.LoginReq) (*proto.LoginRes, error) {
	// start span from context
	span, ctx := ot.StartSpanFromContext(ctx, loginClient)
	u.addSpanTags(span)

	defer span.Finish()
	// send the request to user service
	return u.client.Login(ctx, req)
}

func (u *UserServiceClient) addSpanTags(span ot.Span) {
	span.SetTag(tracing.SpanKind, tracing.SpanKindClient)
	span.SetTag(tracing.Component, tracing.ComponentGrpc)
	span.SetTag(tracing.PeerService, tracing.PeerServiceUserService)
}
