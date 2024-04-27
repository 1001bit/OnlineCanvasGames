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
	flow RunFlow

	clients Children[RoomClient]

	readChan        chan RoomReadMessage
	globalWriteChan chan *message.JSON

	connectedToGame chan struct{}

	id    int
	owner *RoomClient
}

func NewRoomRT() *RoomRT {
	return &RoomRT{
		flow: MakeRunFlow(),

		clients: MakeChildren[RoomClient](),

		readChan:        make(chan RoomReadMessage),
		globalWriteChan: make(chan *message.JSON),

		connectedToGame: make(chan struct{}),

		id:    0,
		owner: nil,
	}
}

func (roomRT *RoomRT) Run(gameRT *GameRT) {
	log.Println("<RoomRT Run>")

	gameRT.rooms.ConnectChild(roomRT)

	defer func() {
		go gameRT.rooms.DisconnectChild(roomRT)
		roomRT.flow.CloseDone()

		log.Println("<RoomRT Done>")
	}()

	// after 5 seconds of start, if there is no client - disconnect the room
	go func() {
		time.Sleep(5 * time.Second)
		if len(roomRT.clients.childMap) == 0 {
			roomRT.flow.Stop()
		}
	}()

	for {
		select {
		case client := <-roomRT.clients.connectChan:
			// When server asked to connect a client
			roomRT.connectClient(client)

			// send rooms json data globally on new client
			go gameRT.requestUpdatingRoomsJSON()

			log.Println("<RoomRT +Client>:", len(roomRT.clients.childMap))

		case client := <-roomRT.clients.disconnectChan:
			// When server asked to disconnect a client
			roomRT.disconnectClient(client)

			// send rooms json data globally on client delete
			go gameRT.requestUpdatingRoomsJSON()

			log.Println("<RoomRT -Client>:", len(roomRT.clients.childMap))

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
	roomRT.clients.childMap[client] = true

	// change owner if no owner yet
	if roomRT.owner == nil {
		roomRT.owner = client
	}
}

// disconnects client from room and sets new owner if owner has left (owner is nil if no clients left, room is deleted after that)
func (roomRT *RoomRT) disconnectClient(client *RoomClient) {
	if _, ok := roomRT.clients.childMap[client]; !ok {
		return
	}

	delete(roomRT.clients.childMap, client)

	// change owner
	if roomRT.owner == client {
		roomRT.owner, _ = roomRT.pickRandomClient()
	}

	// stop room if no clients left after 2 seconds of disconnection
	go func() {
		time.Sleep(2 * time.Second)
		if len(roomRT.clients.childMap) == 0 {
			roomRT.flow.Stop()
		}
	}()
}

// handles message that is read from a client
func (roomRT *RoomRT) handleReadMessage(msg RoomReadMessage) {

}

// write a message to every client
func (roomRT *RoomRT) globalWriteMessage(msg *message.JSON) {
	for client := range roomRT.clients.childMap {
		client.writeChan <- msg
	}
}

// returns random client
func (roomRT *RoomRT) pickRandomClient() (*RoomClient, error) {
	if len(roomRT.clients.childMap) == 0 {
		return nil, ErrNoClients
	}

	k := rand.Intn(len(roomRT.clients.childMap))
	for client := range roomRT.clients.childMap {
		if k == 0 {
			return client, nil
		}
		k--
	}
	return nil, ErrNoClients
}
