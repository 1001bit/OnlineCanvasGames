package server

import (
	"context"
	"database/sql"

	"github.com/1001bit/onlinecanvasgames/services/user/internal/usermodel"
	"github.com/1001bit/onlinecanvasgames/services/user/pkg/userpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer

	userStore *usermodel.UserStore
}

func NewUserServer(userStore *usermodel.UserStore) *UserServer {
	return &UserServer{
		userStore: userStore,
	}
}

func (s *UserServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.UserResponse, error) {
	user, err := s.userStore.GetByName(ctx, req.Username)

	switch err {
	case nil:
		// continue
	case context.DeadlineExceeded:
		return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")
	case sql.ErrNoRows:
		return nil, status.Error(codes.NotFound, "not found")
	default:
		return nil, status.Error(codes.Internal, "something went wrong")
	}

	return &userpb.UserResponse{Username: user.Name, Date: user.Date}, nil
}

func (s *UserServer) LoginUser(ctx context.Context, req *userpb.UserInputRequest) (*userpb.UserResponse, error) {
	user, err := s.userStore.GetByNameAndPassword(ctx, req.Username, req.Password)

	switch err {
	case nil:
		// continue
	case context.DeadlineExceeded:
		return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")
	case usermodel.ErrLogin:
		return nil, status.Error(codes.Unauthenticated, "invalid username or password")
	default:
		return nil, status.Error(codes.Internal, "something went wrong")
	}

	return &userpb.UserResponse{Username: user.Name, Date: user.Date}, nil
}

func (s *UserServer) RegisterUser(ctx context.Context, req *userpb.UserInputRequest) (*userpb.UserResponse, error) {
	user, err := s.userStore.Insert(ctx, req.Username, req.Password)

	switch err {
	case nil:
		// continue
	case context.DeadlineExceeded:
		return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")
	case usermodel.ErrRegister:
		return nil, status.Error(codes.AlreadyExists, "user already exists")
	default:
		return nil, status.Error(codes.Internal, "something went wrong")
	}

	return &userpb.UserResponse{Username: user.Name, Date: user.Date}, nil
}
