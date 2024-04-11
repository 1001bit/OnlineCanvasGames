package sse

import "log"

type GameHub struct {
	clients        map[*Client]bool
	connectChan    chan *Client
	disconnectChan chan *Client

	serverMessageChan chan string

	sse *GamesSSE

	id int
}

func NewGameHub() *GameHub {
	return &GameHub{
		clients:        make(map[*Client]bool),
		connectChan:    make(chan *Client),
		disconnectChan: make(chan *Client),

		serverMessageChan: make(chan string),

		sse: nil,

		id: 0,
	}
}

func (hub *GameHub) Run() {
	log.Println("<GameHub Run>")

	defer func() {
		hub.sse.removeHubIDChan <- hub.id
	}()

	for {
		select {
		case client := <-hub.connectChan:
			hub.clients[client] = true
			client.hub = hub

			log.Println("<GameHub Client Connect>")

		case client := <-hub.disconnectChan:
			delete(hub.clients, client)

			log.Println("<GameHub Client Disconnect>")

		case message := <-hub.serverMessageChan:
			hub.handleServerMessage(message)
		}
	}
}

func (hub *GameHub) handleServerMessage(message string) {
	log.Println("<GameHub Server Message>:", message)
}
