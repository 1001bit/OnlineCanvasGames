package roomnode

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/children"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/runflow"
)

const roomStopWait = 5 * time.Second

var ErrNoClients = errors.New("no clients in the room")

// Struct that contains message and a client who was the message read from
type RoomReadMessage struct {
	client  *RoomClient
	message *message.JSON
}

// Layer of RT which is responsible for handling WS clients
type RoomNode struct {
	Flow runflow.RunFlow

	Clients children.Children[RoomClient]

	readChan        chan RoomReadMessage
	globalWriteChan chan *message.JSON

	connectedToGame chan struct{}

	id    int
	owner *RoomClient
}

func NewRoomNode() *RoomNode {
	return &RoomNode{
		Flow: runflow.MakeRunFlow(),

		Clients: children.MakeChildren[RoomClient](),

		readChan:        make(chan RoomReadMessage),
		globalWriteChan: make(chan *message.JSON),

		connectedToGame: make(chan struct{}),

		id:    0,
		owner: nil,
	}
}

// TODO: Try to get rid of this
func (roomNode *RoomNode) Run(requestUpdatingRoomsJSON func()) {
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
			go requestUpdatingRoomsJSON()

			log.Println("<RoomNode +Client>:", len(roomNode.Clients.ChildMap))

		case client := <-roomNode.Clients.ToDisconnect():
			// When server asked to disconnect a client
			roomNode.disconnectClient(client, stopTimer)

			// Request updaing GameNode's RoomsJSON
			go requestUpdatingRoomsJSON()

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

// connects client to room and makes it owner if no owner exists
func (roomNode *RoomNode) connectClient(client *RoomClient) {
	roomNode.Clients.ChildMap[client] = true

	// change owner if no owner yet
	if roomNode.owner == nil {
		roomNode.owner = client
	}
}

// disconnects client from room and sets new owner if owner has left (owner is nil if no clients left, room is deleted after that)
func (roomNode *RoomNode) disconnectClient(client *RoomClient, stopTimer *time.Timer) {
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
func (roomNode *RoomNode) handleReadMessage(msg RoomReadMessage) {

}

// write a message to every client
func (roomNode *RoomNode) globalWriteMessage(msg *message.JSON) {
	for client := range roomNode.Clients.ChildMap {
		client.writeChan <- msg
	}
}

// returns random client
func (roomNode *RoomNode) pickRandomClient() (*RoomClient, error) {
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
