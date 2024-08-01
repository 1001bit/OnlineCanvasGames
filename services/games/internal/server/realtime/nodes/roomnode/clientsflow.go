package roomnode

import (
	"log"
	"math/rand"
	"time"

	"github.com/1001bit/onlinecanvasgames/services/games/internal/server/realtime/nodes/roomclient"
)

type GameNodeRequester interface {
	RequestUpdatingRoomsJSON()
}

func (roomNode *RoomNode) clientsFlow(requester GameNodeRequester, stopTimer *time.Timer) {
	log.Println("--<RoomNode clientsFlow>")
	defer log.Println("--<RoomNode clientsFlow Done>")

	for {
		select {
		case client := <-roomNode.Clients.ToConnect():
			if roomNode.Clients.ChildrenMap.Length() >= roomNode.gamelogic.GetMaxClients() {
				client.WriteCloseMessage("There are too many players in the room")
				continue
			}

			// When server asked to connect a client
			roomNode.connectClient(client)

			// Request updaing GameNode's RoomsJSON
			go requester.RequestUpdatingRoomsJSON()

			// Notify the gamelogic about new client
			roomNode.gamelogic.JoinClient(client.GetUser().Name, roomNode)

			log.Println("<RoomNode +Client>:", roomNode.Clients.ChildrenMap.Length())

		case client := <-roomNode.Clients.ToDisconnect():
			// When server asked to disconnect a client
			roomNode.disconnectClient(client, stopTimer)

			// Request updaing GameNode's RoomsJSON
			go requester.RequestUpdatingRoomsJSON()

			// Notify the gamelogic about client disconnect
			roomNode.gamelogic.DeleteClient(client.GetUser().Name, roomNode)

			log.Println("<RoomNode -Client>:", roomNode.Clients.ChildrenMap.Length())

		case <-roomNode.Flow.Done():
			return
		}
	}
}

// connects client to room and makes it owner if no owner exists
func (roomNode *RoomNode) connectClient(client *roomclient.RoomClient) {
	roomNode.Clients.ChildrenMap.Set(client.GetUser().Name, client)

	// change owner if no owner yet
	if roomNode.owner == nil {
		roomNode.owner = client
	}
}

// disconnects client from room and sets new owner if owner has left (owner is nil if no clients left, room is deleted after that)
func (roomNode *RoomNode) disconnectClient(client *roomclient.RoomClient, stopTimer *time.Timer) {
	if currentClient, _ := roomNode.Clients.ChildrenMap.Get(client.GetUser().Name); currentClient != client {
		return
	}

	roomNode.Clients.ChildrenMap.Delete(client.GetUser().Name)

	// change owner
	if roomNode.owner == client {
		roomNode.owner, _ = roomNode.pickRandomClient()
	}

	// stop room if no clients left after 2 seconds of disconnection
	stopTimer.Stop()
	stopTimer.Reset(roomStopWait)
}

// returns random client
func (roomNode *RoomNode) pickRandomClient() (*roomclient.RoomClient, error) {
	if roomNode.Clients.ChildrenMap.Length() == 0 {
		return nil, ErrNoClients
	}

	k := rand.Intn(roomNode.Clients.ChildrenMap.Length())

	idMap, rUnlockFunc := roomNode.Clients.ChildrenMap.GetMapForRead()
	defer rUnlockFunc()

	for _, client := range idMap {
		if k == 0 {
			return client, nil
		}
		k--
	}

	return nil, ErrNoClients
}
