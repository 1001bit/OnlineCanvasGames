package sse

import "log"

type GameHub struct {
	clients              map[*Client]bool
	connectClientChan    chan *Client
	disconnectClientChan chan *Client

	globalWriteChan chan string

	sse *GamesSSE

	id int
}

func NewGameHub() *GameHub {
	return &GameHub{
		clients:              make(map[*Client]bool),
		connectClientChan:    make(chan *Client),
		disconnectClientChan: make(chan *Client),

		globalWriteChan: make(chan string),

		sse: nil,

		id: 0,
	}
}

func (hub *GameHub) Run() {
	log.Println("<GameHub Run>")

	defer func() {
		hub.sse.connectHubIDChan <- hub.id
		log.Println("<GameHub Run End>")
	}()

	for {
		select {
		case client := <-hub.connectClientChan:
			hub.connectClient(client)
			log.Println("<GameHub Client Connect>")

		case client := <-hub.disconnectClientChan:
			hub.disconnectClient(client)
			log.Println("<GameHub Client Disconnect>")

		case message := <-hub.globalWriteChan:
			hub.handleGlobalWriteMessage(message)
			log.Println("<GameHub Global Write Message>:", message)
		}
	}
}

func (hub *GameHub) connectClient(client *Client) {
	hub.clients[client] = true
	client.hub = hub
}

func (hub *GameHub) disconnectClient(client *Client) {
	if _, ok := hub.clients[client]; !ok {
		return
	}

	delete(hub.clients, client)
	close(client.writeChan)
}

func (hub *GameHub) handleGlobalWriteMessage(message string) {
	for client := range hub.clients {
		select {
		case client.writeChan <- message:
		default:
			hub.disconnectClientChan <- client
		}
	}
}
