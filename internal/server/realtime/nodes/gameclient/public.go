package gameclient

import (
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
)

func (client *GameClient) WriteMessage(msg *message.JSON) {
	select {
	case client.writeChan <- msg:
		// write message to chan
	default:
		client.Flow.Stop()
	}
}
