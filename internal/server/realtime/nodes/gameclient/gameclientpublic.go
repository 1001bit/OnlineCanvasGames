package gameclient

import (
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
)

func (client *GameClient) WriteMessage(msg *message.JSON) {
	client.writeChan <- msg
}
