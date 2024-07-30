package roomnode

import (
	"log"
	"math/rand"
	"time"

	"github.com/neinBit/ocg-games-service/internal/server/realtime/nodes/roomclient"
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
			if roomNode.Clients.IDMap.Length() >= roomNode.gamelogic.GetMaxClients() {
				client.WriteCloseMessage("There are too many players in the room")
				continue
			}

			// When server asked to connect a client
			roomNode.connectClient(client)

			// Request updaing GameNode's RoomsJSON
			go requester.RequestUpdatingRoomsJSON()

			// Notify the gamelogic about new client
			roomNode.gamelogic.JoinClient(client.GetUser().ID, roomNode)

			log.Println("<RoomNode +Client>:", roomNode.Clients.IDMap.Length())

		case client := <-roomNode.Clients.ToDisconnect():
			// When server asked to disconnect a client
			roomNode.disconnectClient(client, stopTimer)

			// Request updaing GameNode's RoomsJSON
			go requester.RequestUpdatingRoomsJSON()

			// Notify the gamelogic about client disconnect
			roomNode.gamelogic.DeleteClient(client.GetUser().ID, roomNode)

			log.Println("<RoomNode -Client>:", roomNode.Clients.IDMap.Length())

		case <-roomNode.Flow.Done():
			return
		}
	}
}

// connects client to room and makes it owner if no owner exists
func (roomNode *RoomNode) connectClient(client *roomclient.RoomClient) {
	roomNode.Clients.IDMap.Set(client.GetUser().ID, client)

	// change owner if no owner yet
	if roomNode.owner == nil {
		roomNode.owner = client
	}
}

// disconnects client from room and sets new owner if owner has left (owner is nil if no clients left, room is deleted after that)
func (roomNode *RoomNode) disconnectClient(client *roomclient.RoomClient, stopTimer *time.Timer) {
	if currentClient, _ := roomNode.Clients.IDMap.Get(client.GetUser().ID); currentClient != client {
		return
	}

	roomNode.Clients.IDMap.Delete(client.GetUser().ID)

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
	if roomNode.Clients.IDMap.Length() == 0 {
		return nil, ErrNoClients
	}

	k := rand.Intn(roomNode.Clients.IDMap.Length())

	idMap, rUnlockFunc := roomNode.Clients.IDMap.GetMapForRead()
	defer rUnlockFunc()

	for _, client := range idMap {
		if k == 0 {
			return client, nil
		}
		k--
	}

	return nil, ErrNoClients
}
