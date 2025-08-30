package main

import (
	"cmp"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	if err := (func() error {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		addr := ":" + cmp.Or(os.Getenv("GRPC_SERVER_PORT"), "50051")
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		log.Printf("listening on port %s", addr)

		server := grpc.NewServer()

		go func() {
			<-sigCh
			log.Println("shutting down gRPC server...")
			server.GracefulStop()
		}()

		if err = server.Serve(lis); err != nil {
			log.Fatalf("failed to start gRPC server: %v", err)
		}

		return nil
	})(); err != nil {
		log.Fatalf("unknown error occurred: %v", err)
	}
}
