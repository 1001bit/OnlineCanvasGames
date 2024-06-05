package roomnode

import (
	"time"

	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
	rtclient "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/client"
)

func (roomNode *RoomNode) SetRandomID() {
	roomNode.id = int(time.Now().UnixMicro())
}

func (roomNode *RoomNode) ConnectedToGame() <-chan struct{} {
	return roomNode.connectedToGameChan
}

func (roomNode *RoomNode) ConfirmConnectToGame() {
	close(roomNode.connectedToGameChan)
}

func (roomNode *RoomNode) ReadMessage(msg rtclient.MessageWithClient) {
	roomNode.roomplay.HandleReadMessage(msg, roomNode)
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

func (roomNode *RoomNode) GetID() int {
	return roomNode.id
}

func (roomNode *RoomNode) GetPlayersLimit() int {
	return roomNode.roomplay.GetMaxClients()
}

func (roomNode *RoomNode) GetOwnerName() string {
	switch roomNode.owner {
	case nil:
		return "nobody"
	default:
		return roomNode.owner.GetUser().Name
	}
}
