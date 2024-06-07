package roomnode

import (
	"log"
	"math/rand"
	"time"

	rterror "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/error"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/roomclient"
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
			if len(roomNode.Clients.IDMap) >= roomNode.gamelogic.GetMaxClients() {
				client.WriteCloseMessage("There are too many players in the room")
				continue
			}

			// When server asked to connect a client
			roomNode.connectClient(client)

			// Request updaing GameNode's RoomsJSON
			go requester.RequestUpdatingRoomsJSON()

			// Notify the gamelogic about new client
			go roomNode.gamelogic.JoinClient(client.GetUser().ID, roomNode)

			log.Println("<RoomNode +Client>:", len(roomNode.Clients.IDMap))

		case client := <-roomNode.Clients.ToDisconnect():
			// When server asked to disconnect a client
			roomNode.disconnectClient(client, stopTimer)

			// Request updaing GameNode's RoomsJSON
			go requester.RequestUpdatingRoomsJSON()

			log.Println("<RoomNode -Client>:", len(roomNode.Clients.IDMap))

		case <-roomNode.Flow.Done():
			return
		}
	}
}

// connects client to room and makes it owner if no owner exists
func (roomNode *RoomNode) connectClient(client *roomclient.RoomClient) {
	roomNode.Clients.IDMap[client.GetUser().ID] = client

	// change owner if no owner yet
	if roomNode.owner == nil {
		roomNode.owner = client
	}
}

// disconnects client from room and sets new owner if owner has left (owner is nil if no clients left, room is deleted after that)
func (roomNode *RoomNode) disconnectClient(client *roomclient.RoomClient, stopTimer *time.Timer) {
	if roomNode.Clients.IDMap[client.GetUser().ID] != client {
		return
	}

	delete(roomNode.Clients.IDMap, client.GetUser().ID)

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
	if len(roomNode.Clients.IDMap) == 0 {
		return nil, rterror.ErrNoClients
	}

	k := rand.Intn(len(roomNode.Clients.IDMap))
	for _, client := range roomNode.Clients.IDMap {
		if k == 0 {
			return client, nil
		}
		k--
	}
	return nil, rterror.ErrNoClients
}
