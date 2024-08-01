package gameclient

import "github.com/1001bit/onlinecanvasgames/services/games/internal/server/message"

func (client *GameClient) WriteMessage(msg *message.JSON) {
	select {
	case client.writeChan <- msg:
		// write message to chan
	case <-client.Flow.Done():
		// client is done
	}
}
