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
type RoomRT struct {
	Flow runflow.RunFlow

	Clients children.Children[RoomClient]

	readChan        chan RoomReadMessage
	globalWriteChan chan *message.JSON

	connectedToGame chan struct{}

	id    int
	owner *RoomClient
}

func NewRoomRT() *RoomRT {
	return &RoomRT{
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
func (roomRT *RoomRT) Run(requestUpdatingRoomsJSON func()) {
	log.Println("<RoomRT Run>")

	defer func() {
		roomRT.Flow.CloseDone()

		log.Println("<RoomRT Done>")
	}()

	stopTimer := time.NewTimer(roomStopWait)

	for {
		select {
		case client := <-roomRT.Clients.ToConnect():
			// When server asked to connect a client
			roomRT.connectClient(client)

			// Request updaing GameRT's RoomsJSON
			go requestUpdatingRoomsJSON()

			log.Println("<RoomRT +Client>:", len(roomRT.Clients.ChildMap))

		case client := <-roomRT.Clients.ToDisconnect():
			// When server asked to disconnect a client
			roomRT.disconnectClient(client, stopTimer)

			// Request updaing GameRT's RoomsJSON
			go requestUpdatingRoomsJSON()

			log.Println("<RoomRT -Client>:", len(roomRT.Clients.ChildMap))

		case msg := <-roomRT.readChan:
			// Handle messages that were read by all the clients of the room
			roomRT.handleReadMessage(msg)
			log.Println("<RoomRT Read Message>:", msg)

		case msg := <-roomRT.globalWriteChan:
			// Write message to every client if server told to do so
			roomRT.globalWriteMessage(msg)
			log.Println("<RoomRT Global Message>:", msg)

		case <-stopTimer.C:
			// If timer is over, check for clients
			roomRT.stopIfNoClients()

		case <-roomRT.Flow.Stopped():
			// When server asked to stop running
			return
		}
	}
}

func (roomRT *RoomRT) GetID() int {
	return roomRT.id
}

func (roomRT *RoomRT) SetRandomID() {
	roomRT.id = int(time.Now().UnixMicro())
}

func (roomRT *RoomRT) GetOwnerName() string {
	switch roomRT.owner {
	case nil:
		return "nobody"
	default:
		return roomRT.owner.user.Name
	}
}

func (roomRT *RoomRT) ConnectedToGame() <-chan struct{} {
	return roomRT.connectedToGame
}

func (roomRT *RoomRT) ConfirmConnectToGame() {
	close(roomRT.connectedToGame)
}

// connects client to room and makes it owner if no owner exists
func (roomRT *RoomRT) connectClient(client *RoomClient) {
	roomRT.Clients.ChildMap[client] = true

	// change owner if no owner yet
	if roomRT.owner == nil {
		roomRT.owner = client
	}
}

// disconnects client from room and sets new owner if owner has left (owner is nil if no clients left, room is deleted after that)
func (roomRT *RoomRT) disconnectClient(client *RoomClient, stopTimer *time.Timer) {
	if _, ok := roomRT.Clients.ChildMap[client]; !ok {
		return
	}

	delete(roomRT.Clients.ChildMap, client)

	// change owner
	if roomRT.owner == client {
		roomRT.owner, _ = roomRT.pickRandomClient()
	}

	// stop room if no clients left after 2 seconds of disconnection
	stopTimer.Stop()
	stopTimer.Reset(roomStopWait)
}

// handles message that is read from a client
func (roomRT *RoomRT) handleReadMessage(msg RoomReadMessage) {

}

// write a message to every client
func (roomRT *RoomRT) globalWriteMessage(msg *message.JSON) {
	for client := range roomRT.Clients.ChildMap {
		client.writeChan <- msg
	}
}

// returns random client
func (roomRT *RoomRT) pickRandomClient() (*RoomClient, error) {
	if len(roomRT.Clients.ChildMap) == 0 {
		return nil, ErrNoClients
	}

	k := rand.Intn(len(roomRT.Clients.ChildMap))
	for client := range roomRT.Clients.ChildMap {
		if k == 0 {
			return client, nil
		}
		k--
	}
	return nil, ErrNoClients
}

// stops the room if no clients left
func (roomRT *RoomRT) stopIfNoClients() {
	if len(roomRT.Clients.ChildMap) == 0 {
		go roomRT.Flow.Stop()
	}
}
