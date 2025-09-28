package service

import (
	"context"
	pb "service1/pkg/proto"
)

type Service1Grpc struct {
	pb.UnimplementedHelloServiceServer
}

func NewHelloService() *Service1Grpc {
	return &Service1Grpc{}
}

func (s *Service1Grpc) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Reply: "Hello " + req.Name,
	}, nil
}
