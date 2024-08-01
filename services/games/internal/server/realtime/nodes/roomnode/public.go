package roomnode

import (
	"time"

	"github.com/1001bit/onlinecanvasgames/services/games/internal/server/message"
	rtclient "github.com/1001bit/onlinecanvasgames/services/games/internal/server/realtime/client"
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
	roomNode.gamelogic.HandleReadMessage(msg, roomNode)
}

// write a message to every client
func (roomNode *RoomNode) GlobalWriteMessage(msg *message.JSON) {
	idMap, rUnlockFunc := roomNode.Clients.ChildrenMap.GetMapForRead()
	defer rUnlockFunc()

	for _, client := range idMap {
		go client.WriteMessage(msg)
	}
}

// write a message to a single client
func (roomNode *RoomNode) WriteMessageTo(msg *message.JSON, name string) {
	client, ok := roomNode.Clients.ChildrenMap.Get(name)
	if !ok {
		return
	}
	client.WriteMessage(msg)
}

func (roomNode *RoomNode) GetID() int {
	return roomNode.id
}

func (roomNode *RoomNode) GetPlayersLimit() int {
	return roomNode.gamelogic.GetMaxClients()
}

func (roomNode *RoomNode) GetOwnerName() string {
	switch roomNode.owner {
	case nil:
		return "nobody"
	default:
		return roomNode.owner.GetUser().Name
	}
}
