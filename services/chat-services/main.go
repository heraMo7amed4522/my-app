package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "chat-services/proto"
	"chat-services/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "50054" // Default port for chat service
	}

	// Create listener
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// Create gRPC server
	s := grpc.NewServer()

	// Register both regular and streaming servers
	chatServer := server.NewChatServer()
	streamServer := server.NewChatStreamServer()

	// Register services
	pb.RegisterChatServiceServer(s, chatServer)  // For non-streaming methods
	pb.RegisterChatServiceServer(s, streamServer) // For streaming methods

	// Enable reflection for development
	reflection.Register(s)

	log.Printf("Chat Service starting on port %s", port)
	log.Printf("gRPC server listening at %v", lis.Addr())

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down gRPC server...")
		s.GracefulStop()
	}()

	// Start the server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}