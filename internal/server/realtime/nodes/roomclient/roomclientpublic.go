package roomclient

import (
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/rtclient"
)

func (client *RoomClient) WriteMessage(msg *message.JSON) {
	client.writeChan <- msg
}

func (client *RoomClient) GetUser() rtclient.User {
	return client.user
}
