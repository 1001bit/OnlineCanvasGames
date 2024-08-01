package basenode

import (
	"log"

	"github.com/1001bit/onlinecanvasgames/services/games/internal/server/realtime/nodes/roomclient"
)

func (baseNode *BaseNode) clientsFlow() {
	log.Println("--<BaseNode clientsFlow>")
	defer log.Println("--<BaseNode clientsFlow Done>")

	for {
		select {
		case client := <-baseNode.roomsClients.ToConnect():
			// When new WS connection needs to be created
			baseNode.protectRoomClient(client)

		case client := <-baseNode.roomsClients.ToDisconnect():
			// When client is done
			baseNode.deleteRoomClient(client)

		case <-baseNode.Flow.Done():
			return
		}
	}
}

// if there is already client with such ID - stop them. Put a new one
func (baseNode *BaseNode) protectRoomClient(client *roomclient.RoomClient) {
	oldClient, ok := baseNode.roomsClients.ChildrenMap.Get(client.GetUser().Name)
	if ok {
		go oldClient.WriteCloseMessage("This user has just joined another room")
	}
	baseNode.roomsClients.ChildrenMap.Set(client.GetUser().Name, client)
}

// if there is already client with such ID - stop them. Put a new one
func (baseNode *BaseNode) deleteRoomClient(client *roomclient.RoomClient) {
	// can delete only exact client, not just with the same id
	if currentClient, _ := baseNode.roomsClients.ChildrenMap.Get(client.GetUser().Name); currentClient != client {
		return
	}

	baseNode.roomsClients.ChildrenMap.Delete(client.GetUser().Name)
}
