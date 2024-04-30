package roomnode

import (
	"log"
	"math/rand"
	"time"

	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/children"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/roomclient"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/rtclient"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/runflow"
)

const roomStopWait = 5 * time.Second

type GameNodeRequester interface {
	RequestUpdatingRoomsJSON()
}

// Layer of RT which is responsible for handling WS clients
type RoomNode struct {
	Flow runflow.RunFlow

	Clients children.Children[roomclient.RoomClient]

	readChan        chan rtclient.MessageWithClient
	globalWriteChan chan *message.JSON

	connectedToGame chan struct{}

	id    int
	owner *roomclient.RoomClient
}

func NewRoomNode() *RoomNode {
	return &RoomNode{
		Flow: runflow.MakeRunFlow(),

		Clients: children.MakeChildren[roomclient.RoomClient](),

		readChan:        make(chan rtclient.MessageWithClient),
		globalWriteChan: make(chan *message.JSON),

		connectedToGame: make(chan struct{}),

		id:    0,
		owner: nil,
	}
}

func (roomNode *RoomNode) Run(requester GameNodeRequester) {
	log.Println("<RoomNode Run>")

	defer func() {
		roomNode.Flow.CloseDone()

		log.Println("<RoomNode Done>")
	}()

	stopTimer := time.NewTimer(roomStopWait)

	for {
		select {
		case client := <-roomNode.Clients.ToConnect():
			// When server asked to connect a client
			roomNode.connectClient(client)

			// Request updaing GameNode's RoomsJSON
			go requester.RequestUpdatingRoomsJSON()

			log.Println("<RoomNode +Client>:", len(roomNode.Clients.ChildMap))

		case client := <-roomNode.Clients.ToDisconnect():
			// When server asked to disconnect a client
			roomNode.disconnectClient(client, stopTimer)

			// Request updaing GameNode's RoomsJSON
			go requester.RequestUpdatingRoomsJSON()

			log.Println("<RoomNode -Client>:", len(roomNode.Clients.ChildMap))

		case msg := <-roomNode.readChan:
			// Handle messages that were read by all the clients of the room
			roomNode.handleReadMessage(msg)
			log.Println("<RoomNode Read Message>:", msg)

		case msg := <-roomNode.globalWriteChan:
			// Write message to every client if server told to do so
			roomNode.globalWriteMessage(msg)
			log.Println("<RoomNode Global Message>:", msg)

		case <-stopTimer.C:
			// If timer is over, check for clients
			roomNode.stopIfNoClients()

		case <-roomNode.Flow.Stopped():
			// When server asked to stop running
			return
		}
	}
}

// connects client to room and makes it owner if no owner exists
func (roomNode *RoomNode) connectClient(client *roomclient.RoomClient) {
	roomNode.Clients.ChildMap[client] = true

	// change owner if no owner yet
	if roomNode.owner == nil {
		roomNode.owner = client
	}
}

// disconnects client from room and sets new owner if owner has left (owner is nil if no clients left, room is deleted after that)
func (roomNode *RoomNode) disconnectClient(client *roomclient.RoomClient, stopTimer *time.Timer) {
	if _, ok := roomNode.Clients.ChildMap[client]; !ok {
		return
	}

	delete(roomNode.Clients.ChildMap, client)

	// change owner
	if roomNode.owner == client {
		roomNode.owner, _ = roomNode.pickRandomClient()
	}

	// stop room if no clients left after 2 seconds of disconnection
	stopTimer.Stop()
	stopTimer.Reset(roomStopWait)
}

// handles message that is read from a client
func (roomNode *RoomNode) handleReadMessage(msg rtclient.MessageWithClient) {

}

// write a message to every client
func (roomNode *RoomNode) globalWriteMessage(msg *message.JSON) {
	for client := range roomNode.Clients.ChildMap {
		go client.WriteMessage(msg)
	}
}

// returns random client
func (roomNode *RoomNode) pickRandomClient() (*roomclient.RoomClient, error) {
	if len(roomNode.Clients.ChildMap) == 0 {
		return nil, ErrNoClients
	}

	k := rand.Intn(len(roomNode.Clients.ChildMap))
	for client := range roomNode.Clients.ChildMap {
		if k == 0 {
			return client, nil
		}
		k--
	}
	return nil, ErrNoClients
}

// stops the room if no clients left
func (roomNode *RoomNode) stopIfNoClients() {
	if len(roomNode.Clients.ChildMap) == 0 {
		go roomNode.Flow.Stop()
	}
}
