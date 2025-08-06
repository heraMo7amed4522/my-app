package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"chat-services/internal/config"
	"chat-services/internal/repository/postgres"
	"chat-services/internal/server"
	"chat-services/internal/service"
	"chat-services/pkg/database"
	"chat-services/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewPostgresConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repository
	chatRepo := postgres.NewChatRepository(db)

	// Initialize service
	chatService := service.NewChatService(chatRepo)

	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	chatServer := server.NewChatServer(chatService)

	// Register services
	proto.RegisterChatServiceServer(grpcServer, chatServer.UnimplementedChatServiceServer)
	reflection.Register(grpcServer)

	// Start server
	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Chat service starting on port %s", cfg.Port)

	// Graceful shutdown
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down gracefully...")
	grpcServer.GracefulStop()
}
