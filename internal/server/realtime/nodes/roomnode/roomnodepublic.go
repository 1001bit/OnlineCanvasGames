package roomnode

import (
	"time"

	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/rtclient"
)

func (roomNode *RoomNode) SetRandomID() {
	roomNode.id = int(time.Now().UnixMicro())
}

func (roomNode *RoomNode) GetID() int {
	return roomNode.id
}

func (roomNode *RoomNode) GetOwnerName() string {
	switch roomNode.owner {
	case nil:
		return "nobody"
	default:
		return roomNode.owner.GetUser().Name
	}
}

func (roomNode *RoomNode) ConnectedToGame() <-chan struct{} {
	return roomNode.connectedToGameChan
}

func (roomNode *RoomNode) ConfirmConnectToGame() {
	close(roomNode.connectedToGameChan)
}

func (roomNode *RoomNode) ReadMessage(message rtclient.MessageWithClient) {
	// TODO: Handle read message here
}

// write a message to every client
func (roomNode *RoomNode) GlobalWriteMessage(msg *message.JSON) {
	for _, client := range roomNode.Clients.IDMap {
		go client.WriteMessage(msg)
	}
}

// write a message to a single client
func (roomNode *RoomNode) WriteMessageTo(msg *message.JSON, id int) {
	client, ok := roomNode.Clients.IDMap[id]
	if !ok {
		return
	}
	client.WriteMessage(msg)
}
