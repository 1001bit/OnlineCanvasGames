package roomclient

import (
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
)

func (client *RoomClient) WriteMessage(msg *message.JSON) {
	client.writeChan <- msg
}

func (client *RoomClient) GetUser() RoomClientUser {
	return client.user
}
