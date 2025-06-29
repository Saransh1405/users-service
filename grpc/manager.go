package grpc

import (
	"context"
	"fmt"
	"net"
	"users-service/constants"
	"users-service/logger"
	"users-service/proto"
	"users-service/utils/configs"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	server *grpc.Server
	port   int
}

// NewGRPCServer creates a new gRPC server instance
func NewGRPCServer() *GRPCServer {
	// Create server with auth interceptor
	server := grpc.NewServer()

	// Register the user service
	userService := &UserServiceServer{}
	proto.RegisterUserServiceServer(server, userService)

	// Enable reflection for debugging
	reflection.Register(server)

	return &GRPCServer{
		server: server,
	}
}

// Start starts the gRPC server
func (g *GRPCServer) Start(ctx context.Context) error {
	log := logger.GetLoggerWithoutContext()

	// Get configuration
	applicationConfig, err := configs.Get(constants.ApplicationConfig)
	if err != nil {
		log.With(zap.Error(err)).Error("Failed to get application config")
		return err
	}

	// Get gRPC port (default to 9090 if not configured)
	g.port = applicationConfig.GetInt("grpc.port")
	if g.port == 0 {
		g.port = 9090
	}

	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.port))
	if err != nil {
		log.With(zap.Error(err)).Error("Failed to listen")
		return err
	}

	log.Info(fmt.Sprintf("Starting gRPC server on port : %d", g.port))

	// Start server
	go func() {
		if err := g.server.Serve(lis); err != nil {
			log.With(zap.Error(err)).Error("Failed to serve gRPC")
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()
	log.Info("Shutting down gRPC server...")

	// Graceful shutdown
	g.server.GracefulStop()

	return nil
}
