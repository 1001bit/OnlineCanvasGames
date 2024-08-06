package service

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Conn *grpc.ClientConn
}

func NewGRPCClient(host, port string) (*GRPCClient, error) {
	addr := fmt.Sprintf("%s:%s", host, port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return &GRPCClient{
		Conn: conn,
	}, err
}
