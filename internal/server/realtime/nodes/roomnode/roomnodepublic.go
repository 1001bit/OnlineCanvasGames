package roomnode

import (
	"time"

	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/rtclient"
)

func (roomNode *RoomNode) GetID() int {
	return roomNode.id
}

func (roomNode *RoomNode) SetRandomID() {
	roomNode.id = int(time.Now().UnixMicro())
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
	return roomNode.connectedToGame
}

func (roomNode *RoomNode) ConfirmConnectToGame() {
	close(roomNode.connectedToGame)
}

func (roomNode *RoomNode) ReadMessage(message rtclient.MessageWithClient) {
	roomNode.readChan <- message
}
