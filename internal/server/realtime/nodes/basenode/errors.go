package basenode

import "errors"

var (
	ErrNoGame     = errors.New("game does not exist")
	ErrNoRoom     = errors.New("room does not exist")
	ErrCreateRoom = errors.New("could not create a room")
)
