package main

import (
	"cmp"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/KenUtsunomiya/my-rate-limiter/pb/ratelimit/v1"
	"github.com/KenUtsunomiya/my-rate-limiter/pkg/server"
	"github.com/KenUtsunomiya/my-rate-limiter/pkg/valkey"
)

func main() {
	env := cmp.Or(os.Getenv("ENV"), "prod")

	if err := (func() error {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		client, err := valkey.NewClient()
		if err != nil {
			return err
		}
		defer client.Close()

		addr := ":" + cmp.Or(os.Getenv("GRPC_SERVER_PORT"), "50051")
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			return err
		}

		log.Printf("listening on port %s", addr)

		srv := grpc.NewServer()
		pb.RegisterRateLimiterServer(srv, server.NewDefaultServer(*client))

		if env == "dev" {
			reflection.Register(srv)
		}

		shutdownCh := make(chan struct{})
		go func() {
			<-sigCh
			log.Println("shutting down gRPC server...")
			srv.GracefulStop()
			close(shutdownCh)
		}()

		if err = srv.Serve(lis); err != nil {
			return err
		}

		<-shutdownCh
		return nil
	})(); err != nil {
		log.Fatalf("unknown error occurred: %v", err)
	}
}
