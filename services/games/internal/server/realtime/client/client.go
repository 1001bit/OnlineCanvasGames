package rtclient

import "github.com/1001bit/onlinecanvasgames/services/games/internal/server/message"

type User struct {
	Name string
}

type Client interface {
	GetUser() User
}

type MessageWithClient struct {
	Message *message.JSON
	Client  Client
}
