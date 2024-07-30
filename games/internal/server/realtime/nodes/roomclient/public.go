package roomclient

import (
	"github.com/neinBit/ocg-games-service/internal/server/message"
	rtclient "github.com/neinBit/ocg-games-service/internal/server/realtime/client"
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
