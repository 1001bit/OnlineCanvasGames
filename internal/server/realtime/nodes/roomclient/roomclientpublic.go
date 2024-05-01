package roomclient

import (
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/rtclient"
)

func (client *RoomClient) WriteMessage(msg *message.JSON) {
	select {
	case client.writeChan <- msg:
		// write message to writeChan
	default:
		client.Flow.Stop()
	}
}

func (client *RoomClient) GetUser() rtclient.User {
	return client.user
}

// send message to client and stop client
func (client *RoomClient) StopWithMessage(text string) {
	newMessage := &message.JSON{
		Type: "message",
		Body: text,
	}

	select {
	case client.writeChan <- newMessage:
		// write message to chan
	default:
		// just continue
	}

	go client.Flow.Stop()
}
