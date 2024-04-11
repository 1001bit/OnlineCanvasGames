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

	clientMessageChan chan ClientMessage

	ws *GamesWS

	id    int
	owner *Client
}

func NewGameRoom() *GameRoom {
	return &GameRoom{
		clients:        make(map[*Client]bool),
		connectChan:    make(chan *Client),
		disconnectChan: make(chan *Client),

		clientMessageChan: make(chan ClientMessage),

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
			room.clients[client] = true
			client.room = room

			// change owner if no owner yet
			if room.owner == nil {
				room.owner = client
			}

			log.Println("<GameRoom Client Connect>")

		case client := <-room.disconnectChan:
			delete(room.clients, client)

			// change owner or stop room if no clients left
			if room.owner == client {
				room.owner = room.pickRandomClient()
				if room.owner == nil {
					return
				}
			}

			log.Println("<GameRoom Client Disconnect>")

		case message := <-room.clientMessageChan:
			room.handleClientMessage(message)
		}
	}
}

func (room *GameRoom) handleClientMessage(message ClientMessage) {
	log.Println("<GameRoom Message>:", message)
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
