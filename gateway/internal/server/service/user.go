package service

import (
	"context"
	"fmt"
	"io"
)

type User struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

type UserService struct {
	service *Service
}

func NewUserService(host, port string) (*UserService, error) {
	service, err := NewService(host, port)
	if err != nil {
		return nil, err
	}
	return &UserService{
		service: service,
	}, nil
}

// get user from service by username
func (s *UserService) GetUserByName(ctx context.Context, name string) (*User, error) {
	// send get request to service
	msg, err := s.service.request(ctx, "GET", fmt.Sprintf("/%s", name), nil)
	if err != nil {
		return nil, err
	}

	// message types
	switch msg.Type {
	case "user":
		// user message type, it's ok
		user := mapToUser(msg.Body.(map[string]any))
		if user == nil {
			return nil, ErrBadRequest
		}

		return user, nil

	default:
		// some unhandled message type
		return nil, ErrBadRequest
	}
}

// post user to service
func (s *UserService) PostUser(ctx context.Context, body io.ReadCloser) (*User, string) {
	// forward post request from original body
	msg, err := s.service.request(ctx, "POST", "", body)
	if err != nil {
		return nil, ""
	}

	// message types
	switch msg.Type {
	case "user":
		user := mapToUser(msg.Body.(map[string]any))
		if user == nil {
			return nil, ""
		}
		return user, ""

	case "message":
		// message type, something went wrong
		body, ok := msg.Body.(string)
		if !ok {
			return nil, ""
		}
		return nil, body

	default:
		// some unhandled type
		return nil, ""
	}
}

func mapToUser(m map[string]any) *User {
	user := &User{}
	var ok bool

	user.Name, ok = m["name"].(string)
	if !ok {
		return nil
	}

	user.Date, ok = m["date"].(string)
	if !ok {
		return nil
	}

	return user
}
