package forwarder

import (
	service1 "apigateway/pkg/proto/service1"
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Forwarder struct {
	Service1Client service1.HelloServiceClient
}

//setup clients
func NewForwarder(addr string) (*Forwarder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}
	client := service1.NewHelloServiceClient(conn)
	return &Forwarder{
		Service1Client: client,
	}, nil
}
