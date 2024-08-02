package userservice

import (
	"context"

	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/server/service"
	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/userpb"
)

type UserService struct {
	service *service.GRPC
	client  userpb.UserServiceClient
}

func New(host, port string) (*UserService, error) {
	grpcService, err := service.NewGRPCService(host, port)
	if err != nil {
		return nil, err
	}

	service := &UserService{
		service: grpcService,
		client:  userpb.NewUserServiceClient(grpcService.Conn),
	}

	return service, nil
}

// get user from service by username
func (s *UserService) GetUserByName(ctx context.Context, name string) (*userpb.UserResponse, error) {
	return s.client.GetUser(ctx, &userpb.GetUserRequest{Username: name})
}

// login user
func (s *UserService) LoginUser(ctx context.Context, name string, password string) (*userpb.UserResponse, error) {
	input := &userpb.UserInputRequest{
		Username: name,
		Password: password,
	}
	return s.client.LoginUser(ctx, input)
}

// register user
func (s *UserService) RegisterUser(ctx context.Context, name string, password string) (*userpb.UserResponse, error) {
	input := &userpb.UserInputRequest{
		Username: name,
		Password: password,
	}
	return s.client.RegisterUser(ctx, input)
}
