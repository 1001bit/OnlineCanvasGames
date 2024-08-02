package usermodel

import "errors"

var (
	ErrNoUser     = errors.New("incorrect username or password")
	ErrUserExists = errors.New("user with such name already exists")
)
