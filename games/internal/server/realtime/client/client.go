package rtclient

import "github.com/neinBit/ocg-games-service/internal/server/message"

type User struct {
	ID   int
	Name string
}

type Client interface {
	GetUser() User
}

type MessageWithClient struct {
	Message *message.JSON
	Client  Client
}