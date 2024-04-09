package ws

import (
	"log"
	"math/rand"
)

type GameRoom struct {
	WSLayer

	ws      *GamesWS
	clients map[*Client]bool

	id    int
	owner *Client
}

func NewGameRoom() *GameRoom {
	return &GameRoom{
		WSLayer: MakeWSLayer(),

		ws:      nil,
		clients: make(map[*Client]bool),

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
			if room.owner == nil {
				room.owner = client
			}
			client.room = room

			room.ws.connectChan <- client
			log.Println("<GameRoom Client Connect>")

		case client := <-room.disconnectChan:
			delete(room.clients, client)
			if room.owner == client {
				room.owner = room.pickRandomClient()
				if room.owner == nil {
					return
				}
			}

			room.ws.disconnectChan <- client
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
