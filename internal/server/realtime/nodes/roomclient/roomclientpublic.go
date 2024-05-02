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

// send message to client that is going to stop client
func (client *RoomClient) WriteCloseMessage(text string) {
	newMessage := &message.JSON{
		Type: CloseMsgType,
		Body: text,
	}

	client.WriteMessage(newMessage)
}

func (client *RoomClient) GetUser() rtclient.User {
	return client.user
}
