package usermodel

import "errors"

var (
	ErrLogin    = errors.New("incorrect username or password")
	ErrRegister = errors.New("user with such name already exists")
)
