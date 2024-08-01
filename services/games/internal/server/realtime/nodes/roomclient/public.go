package roomclient

import (
	"github.com/1001bit/onlinecanvasgames/services/games/internal/server/message"
	rtclient "github.com/1001bit/onlinecanvasgames/services/games/internal/server/realtime/client"
)

func (client *RoomClient) WriteMessage(msg *message.JSON) {
	select {
	case client.writeChan <- msg:
		// write message to writeChan
	case <-client.Flow.Done():
		// client is done
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
