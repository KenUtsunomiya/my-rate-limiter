package server

import (
	"context"

	"github.com/KenUtsunomiya/my-rate-limiter/pkg/valkey"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/KenUtsunomiya/my-rate-limiter/pb/ratelimit/v1"
)

type Server struct {
	pb.UnimplementedRateLimiterServer
	vkClient valkey.Client
}

func NewServer(client valkey.Client) *Server {
	return &Server{
		vkClient: client,
	}
}

func (s *Server) Hello(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if err := s.vkClient.Hello(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "cannot connect to valkey")
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) Check(ctx context.Context, req *pb.RateLimitRequest) (*pb.RateLimitResponse, error) {
	return &pb.RateLimitResponse{
		Allowed: true,
		Error:   nil,
	}, nil
}
