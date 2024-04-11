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
	clients        map[*Client]bool
	connectChan    chan *Client
	disconnectChan chan *Client

	readChan        chan ClientMessage
	globalWriteChan chan string

	ws *GamesWS

	id    int
	owner *Client
}

func NewGameRoom() *GameRoom {
	return &GameRoom{
		clients:        make(map[*Client]bool),
		connectChan:    make(chan *Client),
		disconnectChan: make(chan *Client),

		readChan:        make(chan ClientMessage),
		globalWriteChan: make(chan string),

		ws: nil,

		id:    0,
		owner: nil,
	}
}

func (room *GameRoom) Run() {
	log.Println("<GameRoom Run>")

	defer func() {
		room.ws.removeRoomIDChan <- room.id
	}()

	for {
		select {
		case client := <-room.connectChan:
			room.connectClient(client)
			log.Println("<GameRoom Client Connect>")

		case client := <-room.disconnectChan:
			room.disconnectClient(client)
			if room.owner == nil {
				return
			}
			log.Println("<GameRoom Client Disconnect>")

		case message := <-room.readChan:
			room.handleReadMessage(message)
			log.Println("<GameRoom Read Message>:", message)

		case message := <-room.globalWriteChan:
			room.handleGlobalWriteMessage(message)
			log.Println("<GameRoom Global Write Message>:", message)
		}
	}
}

func (room *GameRoom) connectClient(client *Client) {
	room.clients[client] = true
	client.room = room

	// change owner if no owner yet
	if room.owner == nil {
		room.owner = client
	}
}

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
}

func (room *GameRoom) handleReadMessage(message ClientMessage) {
	_ = message
}

func (room *GameRoom) handleGlobalWriteMessage(message string) {
	for client := range room.clients {
		select {
		case client.writeChan <- message:
		default:
			room.disconnectChan <- client
		}
	}
}

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
