package main

import (
	"log"
	"net"

	pb "transaction-services/proto"
	"transaction-services/server"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	transactionServer := server.NewTransactionServer()
	pb.RegisterTransactionServiceServer(s, transactionServer)

	log.Println("Transaction service is running on port 50054...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}