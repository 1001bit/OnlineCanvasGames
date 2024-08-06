package userservice

import (
	"context"

	service "github.com/1001bit/onlinecanvasgames/services/gateway/pkg/client"
	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/userpb"
)

type Client struct {
	service *service.GRPCClient
	client  userpb.UserServiceClient
}

func NewClient(host, port string) (*Client, error) {
	grpcService, err := service.NewGRPCClient(host, port)
	if err != nil {
		return nil, err
	}

	service := &Client{
		service: grpcService,
		client:  userpb.NewUserServiceClient(grpcService.Conn),
	}

	return service, nil
}

// get user from service by username
func (s *Client) GetUserByName(ctx context.Context, name string) (*userpb.UserResponse, error) {
	return s.client.GetUser(ctx, &userpb.GetUserRequest{Username: name})
}

// login user
func (s *Client) LoginUser(ctx context.Context, name string, password string) (*userpb.UserResponse, error) {
	input := &userpb.UserInputRequest{
		Username: name,
		Password: password,
	}
	return s.client.LoginUser(ctx, input)
}

// register user
func (s *Client) RegisterUser(ctx context.Context, name string, password string) (*userpb.UserResponse, error) {
	input := &userpb.UserInputRequest{
		Username: name,
		Password: password,
	}
	return s.client.RegisterUser(ctx, input)
}
