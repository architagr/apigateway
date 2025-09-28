package main

import (
	"log"
	"net"

	pb "service1/pkg/proto"
	"service1/pkg/service"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("error starting the server: %v", err)
	}

	grpcServer := grpc.NewServer()

	serviceImpl := service.NewHelloService()
	pb.RegisterHelloServiceServer(grpcServer, serviceImpl)

	log.Println("gRPC server running on :5001")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
