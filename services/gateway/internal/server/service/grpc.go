package service

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPC struct {
	Conn *grpc.ClientConn
}

func NewGRPCService(host, port string) (*GRPC, error) {
	addr := fmt.Sprintf("%s:%s", host, port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return &GRPC{
		Conn: conn,
	}, err
}
