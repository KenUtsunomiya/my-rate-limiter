package server

import (
	"cmp"
	"context"
	"log"
	"os"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pb "github.com/KenUtsunomiya/my-rate-limiter/pb/ratelimit/v1"
	"github.com/KenUtsunomiya/my-rate-limiter/pkg/ratelimit"
	"github.com/KenUtsunomiya/my-rate-limiter/pkg/valkey"
)

type Server struct {
	pb.UnimplementedRateLimiterServer
	vkClient valkey.Client
	rl       rateLimiter
}

type rateLimiter interface {
	Allow(ctx context.Context, userId string, method string, resource string) (bool, error)
}

func NewServer(client valkey.Client, rl rateLimiter) *Server {
	return &Server{
		vkClient: client,
		rl:       rl,
	}
}

func NewDefaultServer(client valkey.Client) *Server {
	maxTokens, _ := strconv.ParseFloat(cmp.Or(os.Getenv("RATELIMIT_MAX_TOKENS"), "1000"), 64)
	refillRate, _ := strconv.ParseFloat(cmp.Or(os.Getenv("RATELIMIT_REFILL_RATE"), "1000"), 64)
	return NewServer(
		client,
		ratelimit.NewRateLimiter(client, "ratelimit:", maxTokens, refillRate))
}

func (s *Server) Hello(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if err := s.vkClient.Ping(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "cannot connect to valkey")
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) Check(ctx context.Context, req *pb.RateLimitRequest) (*pb.RateLimitResponse, error) {
	log.Printf("check rate limit: user_id=%s, method=%s, resource=%s", req.UserId, req.Method, req.Resource)

	allowed, err := s.rl.Allow(ctx, req.UserId, req.Method, req.Resource)
	if err != nil {
		return &pb.RateLimitResponse{
			Allowed: &wrapperspb.BoolValue{Value: false},
			Error: &pb.Error{
				Code:    pb.Error_INTERNAL_ERROR,
				Message: err.Error(),
			},
		}, nil
	}

	return &pb.RateLimitResponse{
		Allowed: &wrapperspb.BoolValue{Value: allowed},
		Error:   nil,
	}, nil
}
