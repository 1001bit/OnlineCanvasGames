package roomnode

import "time"

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
		return roomNode.owner.user.Name
	}
}

func (roomNode *RoomNode) ConnectedToGame() <-chan struct{} {
	return roomNode.connectedToGame
}

func (roomNode *RoomNode) ConfirmConnectToGame() {
	close(roomNode.connectedToGame)
}
