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
	client  *RoomRTClient
	message MessageJSON
}

// Layer of RT which is responsible for handling WS clients
type RoomRT struct {
	gameRT *GameRT

	done chan struct{}

	clients              map[*RoomRTClient]bool
	connectClientChan    chan *RoomRTClient
	disconnectClientChan chan *RoomRTClient

	readChan chan RoomReadMessage

	connectedToRT chan struct{}

	id    int
	owner *RoomRTClient
}

func NewRoomRT() *RoomRT {
	return &RoomRT{
		gameRT: nil,

		done: make(chan struct{}),

		clients:              make(map[*RoomRTClient]bool),
		connectClientChan:    make(chan *RoomRTClient),
		disconnectClientChan: make(chan *RoomRTClient),

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

		// send rooms json data globally on room stop
		go func() {
			roomRT.gameRT.globalWriteChan <- roomRT.gameRT.prepareRoomsMessage()
		}()
	}()

	// send rooms json data globally on room start
	go func() {
		roomRT.gameRT.globalWriteChan <- roomRT.gameRT.prepareRoomsMessage()
	}()

	// after 5 seconds of start, if there is no client - disconnect the room
	go func() {
		time.Sleep(5 * time.Second)
		if len(roomRT.clients) == 0 {
			roomRT.gameRT.disconnectRoomChan <- roomRT
		}
	}()

	for {
		select {
		// Clients
		case client := <-roomRT.connectClientChan:
			roomRT.connectClient(client)

			// send rooms json data globally on new client
			go func() {
				roomRT.gameRT.globalWriteChan <- roomRT.gameRT.prepareRoomsMessage()
			}()

			log.Println("<RoomRT +Client>:", len(roomRT.clients))

		case client := <-roomRT.disconnectClientChan:
			roomRT.disconnectClient(client)

			// send rooms json data globally on client delete
			go func() {
				roomRT.gameRT.globalWriteChan <- roomRT.gameRT.prepareRoomsMessage()
			}()

			log.Println("<RoomRT -Client>:", len(roomRT.clients))

		// Handle messages that were read by all the clients of the room
		case message := <-roomRT.readChan:
			roomRT.handleReadMessage(message)
			log.Println("<RoomRT Read Message>:", message)

		// When game closed roomRT.done chan
		case <-roomRT.done:
			return
		}
	}
}

func (roomRT *RoomRT) GetID() int {
	return roomRT.id
}

// connects client to room and makes it owner if no owner exists
func (roomRT *RoomRT) connectClient(client *RoomRTClient) {
	roomRT.clients[client] = true
	client.roomRT = roomRT

	// change owner if no owner yet
	if roomRT.owner == nil {
		roomRT.owner = client
	}

	// add client to RT's list
	roomRT.gameRT.rt.registerRoomClientChan <- client

	go client.Run()
}

// disconnects client from room and sets new owner if owner has left (owner is nil if no clients left, room is deleted after that)
func (roomRT *RoomRT) disconnectClient(client *RoomRTClient) {
	if _, ok := roomRT.clients[client]; !ok {
		return
	}

	delete(roomRT.clients, client)
	close(client.done)

	// change owner
	if roomRT.owner == client {
		roomRT.owner, _ = roomRT.pickRandomClient()
	}

	// delete room if no clients left
	if len(roomRT.clients) == 0 {
		roomRT.gameRT.disconnectRoomChan <- roomRT
	}

	// remove client from RT's list
	roomRT.gameRT.rt.unregisterRoomClientChan <- client
}

// handles message that is read from a client
func (roomRT *RoomRT) handleReadMessage(message RoomReadMessage) {

}

// returns random client
func (roomRT *RoomRT) pickRandomClient() (*RoomRTClient, error) {
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
