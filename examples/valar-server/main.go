package main

// Example main.go for Valar gRPC server
// Shows how to start the gRPC server and register services

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	emailv1 "github.com/Cloud9Money/maia/proto/email/v1"
	smsv1 "github.com/Cloud9Money/maia/proto/sms/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Configuration
	grpcPort := getEnv("GRPC_PORT", "50051")
	httpPort := getEnv("HTTP_PORT", "8080")

	log.Printf("Starting Valar gRPC server on port %s", grpcPort)

	// Create TCP listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(10*1024*1024), // 10MB
		grpc.MaxSendMsgSize(10*1024*1024), // 10MB
	)

	// TODO: Initialize your dependencies
	// resendClient := initializeResendClient()
	// logger := initializeLogger()

	// Register Email Service
	// emailServer := grpcserver.NewEmailServer(resendClient, logger)
	// emailv1.RegisterEmailServiceServer(grpcServer, emailServer)

	// Register SMS Service
	// smsServer := grpcserver.NewSMSServer(smsProvider, logger)
	// smsv1.RegisterSMSServiceServer(grpcServer, smsServer)

	// Enable gRPC reflection for debugging with grpcurl
	reflection.Register(grpcServer)

	// Start gRPC server in goroutine
	go func() {
		log.Printf("✅ gRPC server listening on :%s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// TODO: Start HTTP server for health checks and metrics on httpPort

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("gRPC server stopped")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
