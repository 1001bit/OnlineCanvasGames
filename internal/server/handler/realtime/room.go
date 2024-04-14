package realtime

import (
	"log"
	"math/rand"
)

type RoomRTClientMessage struct {
	client *RoomRTClient
	text   string
}

// Layer of RT which is responsible for handling WS clients
type RoomRT struct {
	gameRT *GameRT

	clients              map[*RoomRTClient]bool
	connectClientChan    chan *RoomRTClient
	disconnectClientChan chan *RoomRTClient

	readChan chan RoomRTClientMessage

	connectedToRT chan struct{}

	id    int
	owner *RoomRTClient
}

func NewRoomRT() *RoomRT {
	return &RoomRT{
		gameRT: nil,

		clients:              make(map[*RoomRTClient]bool),
		connectClientChan:    make(chan *RoomRTClient),
		disconnectClientChan: make(chan *RoomRTClient),

		readChan: make(chan RoomRTClientMessage),

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

	for {
		select {
		case client := <-roomRT.connectClientChan:
			roomRT.connectClient(client)
			log.Println("<RoomRT +Client>:", len(roomRT.clients))

		case client := <-roomRT.disconnectClientChan:
			roomRT.disconnectClient(client)
			if roomRT.owner == nil {
				return
			}
			log.Println("<RoomRT -Client>:", len(roomRT.clients))

		case message := <-roomRT.readChan:
			roomRT.handleReadMessage(message)
			log.Println("<RoomRT Read Message>:", message)
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
	roomRT.gameRT.rt.connectRoomClientChan <- client
}

// disconnects client from room and sets new owner if owner has left (owner is nil if no clients left, room is deleted after that)
func (roomRT *RoomRT) disconnectClient(client *RoomRTClient) {
	if _, ok := roomRT.clients[client]; !ok {
		return
	}

	delete(roomRT.clients, client)
	close(client.writeChan)

	// change owner
	if roomRT.owner == client {
		roomRT.owner = roomRT.pickRandomClient()
	}

	// remove client from RT's list
	roomRT.gameRT.rt.disconnectRoomClientChan <- client
}

// handles message that is read from a client
func (roomRT *RoomRT) handleReadMessage(message RoomRTClientMessage) {

}

// returns random client
func (roomRT *RoomRT) pickRandomClient() *RoomRTClient {
	if len(roomRT.clients) == 0 {
		return nil
	}

	k := rand.Intn(len(roomRT.clients))
	for client := range roomRT.clients {
		if k == 0 {
			return client
		}
		k--
	}
	return nil
}
