package ws

import (
	"log"
	"math/rand"
)

type ClientMessage struct {
	client *Client
	text   string
}

type GameRoom struct {
	clients              map[*Client]bool
	connectClientChan    chan *Client
	disconnectClientChan chan *Client

	// channel of client's read messages
	readChan chan ClientMessage
	// channel of messages that must be send to every client
	globalWriteChan chan string

	ws *GamesWS

	// channel that is closed when GameRoom is connectedToWS to ws
	connectedToWS chan struct{}

	id    int
	owner *Client
}

func NewGameRoom() *GameRoom {
	return &GameRoom{
		clients:              make(map[*Client]bool),
		connectClientChan:    make(chan *Client),
		disconnectClientChan: make(chan *Client),

		readChan:        make(chan ClientMessage),
		globalWriteChan: make(chan string),

		ws: nil,

		connectedToWS: make(chan struct{}),

		id:    0,
		owner: nil,
	}
}

func (room *GameRoom) Run() {
	log.Println("<GameRoom Run>")

	defer func() {
		room.ws.disconnectRoomChan <- room
		log.Println("<GameRoom Run End>")
	}()

	for {
		select {
		case client := <-room.connectClientChan:
			room.connectClient(client)
			log.Println("<GameRoom +Client>:", len(room.clients))

		case client := <-room.disconnectClientChan:
			room.disconnectClient(client)
			if room.owner == nil {
				return
			}
			log.Println("<GameRoom -Client>:", len(room.clients))

		case message := <-room.readChan:
			room.handleReadMessage(message)
			log.Println("<GameRoom Read Message>:", message)

		case message := <-room.globalWriteChan:
			room.handleGlobalWriteMessage(message)
			log.Println("<GameRoom Global Write Message>:", message)
		}
	}
}

func (room *GameRoom) GetID() int {
	return room.id
}

// stops the goroutine until connectedToWS is closed
func (room *GameRoom) waitUntilConnectedToWS() {
	<-room.connectedToWS
}

// closes connected chan
func (room *GameRoom) confirmConnectToWS() {
	close(room.connectedToWS)
}

// connects a client to the room
func (room *GameRoom) connectClient(client *Client) {
	room.clients[client] = true
	client.room = room

	// change owner if no owner yet
	if room.owner == nil {
		room.owner = client
	}

	// add client to ws list
	room.ws.connectClientChan <- client
}

// disconnects a client from the room
func (room *GameRoom) disconnectClient(client *Client) {
	if _, ok := room.clients[client]; !ok {
		return
	}

	delete(room.clients, client)
	close(client.writeChan)

	// change owner
	if room.owner == client {
		room.owner = room.pickRandomClient()
	}

	// remove client from ws list
	room.ws.disconnectClientChan <- client
}

// handles messages that are read from a client
func (room *GameRoom) handleReadMessage(message ClientMessage) {
	_ = message
}

// writes a message to every client
func (room *GameRoom) handleGlobalWriteMessage(message string) {
	for client := range room.clients {
		select {
		// write message to client's writeChan
		case client.writeChan <- message:
		// disconnect client if cannot write the message
		default:
			room.disconnectClientChan <- client
		}
	}
}

// TODO: Make it thread safe
// returns random client of the room
func (room *GameRoom) pickRandomClient() *Client {
	if len(room.clients) == 0 {
		return nil
	}

	k := rand.Intn(len(room.clients))
	for client := range room.clients {
		if k == 0 {
			return client
		}
		k--
	}
	return nil
}
