package basenode

import (
	"log"

	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/roomclient"
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
	oldClient, ok := baseNode.roomsClients.IDMap.Get(client.GetUser().ID)
	if ok {
		go oldClient.WriteCloseMessage("This user has just joined another room")
	}
	baseNode.roomsClients.IDMap.Set(client.GetUser().ID, client)
}

// if there is already client with such ID - stop them. Put a new one
func (baseNode *BaseNode) deleteRoomClient(client *roomclient.RoomClient) {
	// can delete only exact client, not just with the same id
	if currentClient, _ := baseNode.roomsClients.IDMap.Get(client.GetUser().ID); currentClient != client {
		return
	}

	baseNode.roomsClients.IDMap.Delete(client.GetUser().ID)
}
