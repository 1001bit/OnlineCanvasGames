package usermodel

import "errors"

var (
	ErrNoUserExists = errors.New("user with such name doesn't exist")
	ErrNoSuchUser   = errors.New("incorrect username or password")
	ErrUserExists   = errors.New("user with such name already exists")
)
