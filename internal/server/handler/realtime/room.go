package realtime

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
)

var ErrNoClients = errors.New("no clients in the room")

// Struct that contains message and a client who was the message read from
type RoomReadMessage struct {
	client  *RoomClient
	message *message.JSON
}

// Layer of RT which is responsible for handling WS clients
type RoomRT struct {
	gameRT *GameRT

	flow RunFlow

	clients              map[*RoomClient]bool
	connectClientChan    chan *RoomClient
	disconnectClientChan chan *RoomClient

	readChan        chan RoomReadMessage
	globalWriteChan chan *message.JSON

	connectedToGame chan struct{}

	id    int
	owner *RoomClient
}

func NewRoomRT() *RoomRT {
	return &RoomRT{
		gameRT: nil,

		flow: MakeRunFlow(),

		clients:              make(map[*RoomClient]bool),
		connectClientChan:    make(chan *RoomClient),
		disconnectClientChan: make(chan *RoomClient),

		readChan:        make(chan RoomReadMessage),
		globalWriteChan: make(chan *message.JSON),

		connectedToGame: make(chan struct{}),

		id:    0,
		owner: nil,
	}
}

func (roomRT *RoomRT) Run() {
	log.Println("<RoomRT Run>")

	defer func() {
		roomRT.gameRT.disconnectRoomChan <- roomRT
		roomRT.flow.CloseDone()

		log.Println("<RoomRT Done>")
	}()

	// after 5 seconds of start, if there is no client - disconnect the room
	go func() {
		time.Sleep(5 * time.Second)
		if len(roomRT.clients) == 0 {
			roomRT.flow.Stop()
		}
	}()

	for {
		select {
		case client := <-roomRT.connectClientChan:
			// When server asked to connect a client
			roomRT.connectClient(client)

			// send rooms json data globally on new client
			go roomRT.gameRT.requestUpdatingRoomsJSON()

			log.Println("<RoomRT +Client>:", len(roomRT.clients))

		case client := <-roomRT.disconnectClientChan:
			// When server asked to disconnect a client
			roomRT.disconnectClient(client)

			// send rooms json data globally on client delete
			go roomRT.gameRT.requestUpdatingRoomsJSON()

			log.Println("<RoomRT -Client>:", len(roomRT.clients))

		case msg := <-roomRT.readChan:
			// Handle messages that were read by all the clients of the room
			roomRT.handleReadMessage(msg)
			log.Println("<RoomRT Read Message>:", msg)

		case msg := <-roomRT.globalWriteChan:
			// Write message to every client if server told to do so
			roomRT.globalWriteMessage(msg)
			log.Println("<RoomRT Global Message>:", msg)

		case <-roomRT.flow.Stopped():
			// When server asked to stop running
			return
		}
	}
}

func (roomRT *RoomRT) GetID() int {
	return roomRT.id
}

// connects client to room and makes it owner if no owner exists
func (roomRT *RoomRT) connectClient(client *RoomClient) {
	roomRT.clients[client] = true
	client.roomRT = roomRT

	// change owner if no owner yet
	if roomRT.owner == nil {
		roomRT.owner = client
	}

	// add client to RT's list
	roomRT.gameRT.rt.registerRoomClientChan <- client

	// Not running client, as it's being ran right upon it's creation
}

// disconnects client from room and sets new owner if owner has left (owner is nil if no clients left, room is deleted after that)
func (roomRT *RoomRT) disconnectClient(client *RoomClient) {
	if _, ok := roomRT.clients[client]; !ok {
		return
	}

	delete(roomRT.clients, client)

	// change owner
	if roomRT.owner == client {
		roomRT.owner, _ = roomRT.pickRandomClient()
	}

	// stop room if no clients left after 2 seconds of disconnection
	go func() {
		time.Sleep(2 * time.Second)
		if len(roomRT.clients) == 0 {
			roomRT.flow.Stop()
		}
	}()

	// remove client from RT's list
	roomRT.gameRT.rt.unregisterRoomClientChan <- client
}

// handles message that is read from a client
func (roomRT *RoomRT) handleReadMessage(msg RoomReadMessage) {

}

// write a message to every client
func (roomRT *RoomRT) globalWriteMessage(msg *message.JSON) {
	for client := range roomRT.clients {
		client.writeChan <- msg
	}
}

// returns random client
func (roomRT *RoomRT) pickRandomClient() (*RoomClient, error) {
	if len(roomRT.clients) == 0 {
		return nil, ErrNoClients
	}

	k := rand.Intn(len(roomRT.clients))
	for client := range roomRT.clients {
		if k == 0 {
			return client, nil
		}
		k--
	}
	return nil, ErrNoClients
}
