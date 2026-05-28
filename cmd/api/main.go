package main

import (
	"log"
	"net"

	"github.com/Krchnk/go-micro/internal/config"
	usersv1 "github.com/Krchnk/go-micro/internal/gen/users/v1"
	"github.com/Krchnk/go-micro/internal/grpcapi"
	"github.com/Krchnk/go-micro/internal/monitoring"
	"github.com/Krchnk/go-micro/internal/users"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()
	monitoring.StartMetricsServer(cfg.MetricsAddr)
	repo := users.NewInMemoryRepository()
	service := users.NewService(repo)

	listener, err := net.Listen("tcp", cfg.GRPCAddr())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcapi.JWTUnaryInterceptor()),
	)
	usersv1.RegisterUserServiceServer(grpcServer, grpcapi.NewServer(service, cfg))

	log.Printf("user gRPC service is running on %s", cfg.GRPCAddr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
