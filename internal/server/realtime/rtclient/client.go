package rtclient

import "github.com/1001bit/OnlineCanvasGames/internal/server/message"

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
