package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	pb "wallet-services/proto"
	"wallet-services/server"
)

func main() {
	// Create a TCP listener on port 50055
	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Initialize user service client
	userClient, err := server.NewUserServiceClient()
	if err != nil {
		log.Fatalf("Failed to initialize user service client: %v", err)
	}
	defer userClient.Close()

	// Create a new gRPC server with authentication middleware
	s := grpc.NewServer(
		grpc.UnaryInterceptor(server.AuthInterceptor(userClient)),
	)

	// Create wallet server instance
	walletServer := server.NewWalletServer()

	// Register the wallet service
	pb.RegisterWalletServiceServer(s, walletServer)

	log.Println("Wallet service is running on port 50055 with authentication middleware...")

	// Start serving
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}