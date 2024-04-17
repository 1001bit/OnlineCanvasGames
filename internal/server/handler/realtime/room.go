package realtime

import (
	"errors"
	"log"
	"math/rand"
	"time"
)

var ErrNoClients = errors.New("no clients in the room")

// Struct that contains message and a client who was the message read from
type RoomReadMessage struct {
	client  *RoomClient
	message MessageJSON
}

// Layer of RT which is responsible for handling WS clients
type RoomRT struct {
	gameRT *GameRT

	stopChan chan struct{}
	doneChan chan struct{}

	clients              map[*RoomClient]bool
	connectClientChan    chan *RoomClient
	disconnectClientChan chan *RoomClient

	readChan chan RoomReadMessage

	connectedToRT chan struct{}

	id    int
	owner *RoomClient
}

func NewRoomRT() *RoomRT {
	return &RoomRT{
		gameRT: nil,

		stopChan: make(chan struct{}),
		doneChan: make(chan struct{}),

		clients:              make(map[*RoomClient]bool),
		connectClientChan:    make(chan *RoomClient),
		disconnectClientChan: make(chan *RoomClient),

		readChan: make(chan RoomReadMessage),

		connectedToRT: make(chan struct{}),

		id:    0,
		owner: nil,
	}
}

func (roomRT *RoomRT) Run() {
	log.Println("<RoomRT Run>")

	defer func() {
		roomRT.gameRT.disconnectRoomChan <- roomRT
		log.Println("<RoomRT Run End>")
	}()

	// after 5 seconds of start, if there is no client - disconnect the room
	go func() {
		time.Sleep(5 * time.Second)
		if len(roomRT.clients) == 0 {
			roomRT.Stop()
		}
	}()

	for {
		select {
		case client := <-roomRT.connectClientChan:
			// When server asked to connect a client
			roomRT.connectClient(client)

			// send rooms json data globally on new client
			roomRT.gameRT.globallyWriteRoomsMessage()

			log.Println("<RoomRT +Client>:", len(roomRT.clients))

		case client := <-roomRT.disconnectClientChan:
			// When server asked to disconnect a client
			roomRT.disconnectClient(client)

			// send rooms json data globally on client delete
			roomRT.gameRT.globallyWriteRoomsMessage()

			log.Println("<RoomRT -Client>:", len(roomRT.clients))

		case message := <-roomRT.readChan:
			// Handle messages that were read by all the clients of the room
			roomRT.handleReadMessage(message)
			log.Println("<RoomRT Read Message>:", message)

		case <-roomRT.stopChan:
			// When server asked to stop running
			return

		case <-roomRT.doneChan:
			// When parent closed done chan
			return
		}
	}
}

func (roomRT *RoomRT) Stop() {
	roomRT.stopChan <- struct{}{}
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
	close(client.doneChan)

	// change owner
	if roomRT.owner == client {
		roomRT.owner, _ = roomRT.pickRandomClient()
	}

	// stop room if no clients left
	if len(roomRT.clients) == 0 {
		go roomRT.Stop()
	}

	// remove client from RT's list
	roomRT.gameRT.rt.unregisterRoomClientChan <- client
}

// handles message that is read from a client
func (roomRT *RoomRT) handleReadMessage(message RoomReadMessage) {

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
