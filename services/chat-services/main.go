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
	port := os.Getenv("PORT")
	if port == "" {
		port = "50054"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}
	s := grpc.NewServer()
	// Only register the stream server since it embeds ChatServer and handles both streaming and non-streaming methods
	streamServer := server.NewChatStreamServer()
	pb.RegisterChatServiceServer(s, streamServer)
	reflection.Register(s)
	log.Printf("Chat Service starting on port %s", port)
	log.Printf("gRPC server listening at %v", lis.Addr())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Shutting down gRPC server...")
		s.GracefulStop()
	}()
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
