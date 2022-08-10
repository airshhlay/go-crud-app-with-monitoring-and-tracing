package server

import (
	"context"
	"fmt"
	config "itemService/config"
	constants "itemService/constants"
	db "itemService/db"
	customErr "itemService/errors"
	pb "itemService/proto"
	"net"

	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	pb.UnimplementedItemServiceServer
	handler Handler
	logger  *zap.Logger
}

func (s *Server) StartServer(config *config.Config, logger *zap.Logger, dbManager *db.DbManager, redisManager *db.RedisManager) {
	s.handler = Handler{
		config:       config,
		dbManager:    dbManager,
		redisManager: redisManager,
		logger:       logger,
	}
	s.logger = logger

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Port))
	if err != nil {
		s.logger.Fatal(
			constants.ERROR_SERVER_START_FAIL_MSG,
			zap.Error(err),
		)
	}

	server := grpc.NewServer()
	pb.RegisterItemServiceServer(server, s)
	logger.Info(
		constants.INFO_SERVER_START_MSG,
		zap.Any("port", listener.Addr()),
	)

	reflection.Register(server)
	err = server.Serve(listener)
	if err != nil {
		logger.Fatal(
			constants.ERROR_SERVER_START_FAIL_MSG,
			zap.Error(err),
		)
	}
}

func (s *Server) DeleteFav(ctx context.Context, req *pb.DeleteFavReq) (*pb.DeleteFavRes, error) {
	err := s.handler.DeleteFavourite(req.UserId, req.ItemId, req.ShopId)
	if err != nil {
		v, ok := err.(*customErr.Error)
		if !ok {
			s.logger.Error(
				constants.ERROR_TYPECAST_MSG,
				zap.Error(err),
			)
			return &pb.DeleteFavRes{
				ErrorCode: constants.ERROR_TYPECAST,
			}, nil
		}
		return &pb.DeleteFavRes{
			ErrorCode: v.ErrorCode,
		}, nil
	}
	return &pb.DeleteFavRes{
		ErrorCode: -1,
	}, nil
}

func (s *Server) AddFav(ctx context.Context, req *pb.AddFavReq) (*pb.AddFavRes, error) {
	item, err := s.handler.AddItemToUserFavList(req.ItemId, req.ShopId, req.UserId)

	if err != nil {
		v, ok := err.(*customErr.Error)
		if !ok {
			s.logger.Error(
				constants.ERROR_TYPECAST_MSG,
				zap.Error(err),
			)
			return &pb.AddFavRes{
				ErrorCode: constants.ERROR_TYPECAST,
			}, nil
		}
		return &pb.AddFavRes{
			ErrorCode: v.ErrorCode,
		}, nil
	}
	return &pb.AddFavRes{
		ErrorCode: -1,
		Item:      item,
	}, nil
}

func (s *Server) GetFavList(ctx context.Context, req *pb.GetFavListReq) (*pb.GetFavListRes, error) {

	items, totalPages, err := s.handler.GetUserFavourites(req.UserId, req.Page)
	if err != nil {
		v, ok := err.(*customErr.Error)
		if !ok {
			s.logger.Error(
				constants.ERROR_TYPECAST_MSG,
				zap.Error(err),
			)
			return &pb.GetFavListRes{
				ErrorCode: constants.ERROR_TYPECAST,
			}, nil
		}
		return &pb.GetFavListRes{
			ErrorCode: v.ErrorCode,
		}, nil
	}

	return &pb.GetFavListRes{
		ErrorCode:  -1,
		Items:      items,
		TotalPages: totalPages,
	}, nil
}
