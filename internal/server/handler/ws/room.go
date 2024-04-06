package ws

import (
	"log"
)

type GameRoom struct {
	WSLayer
	ws      *GamesWS
	clients map[*Client]bool
}

func NewGameRoom() *GameRoom {
	return &GameRoom{
		WSLayer: MakeWSLayer(),
		clients: make(map[*Client]bool),
	}
}

func (room *GameRoom) Run() {
	log.Println("<GameRoom Run>")

	for {
		select {
		case client := <-room.connectChan:
			room.clients[client] = true
			room.ws.connectChan <- client
			log.Println("<GameRoom Connect>")

		case client := <-room.disconnectChan:
			delete(room.clients, client)
			room.ws.disconnectChan <- client
			log.Println("<GameRoom Disconnect>")

		case message := <-room.messageChan:
			room.handleMessage(message)
		}
	}
}

func (room *GameRoom) handleMessage(message string) {
	log.Println("<GameRoom Message>:", message)
}
