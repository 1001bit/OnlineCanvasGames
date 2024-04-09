package sse

import "log"

type GameHub struct {
	SSELayer

	sse *GamesSSE

	id int
}

func NewGameHub() *GameHub {
	return &GameHub{
		SSELayer: MakeSSELayer(),

		sse: nil,

		id: 0,
	}
}

func (hub *GameHub) Run() {
	log.Println("<GameHub Run>")

	for {
		select {
		case client := <-hub.connect:
			hub.clients[client] = true
			hub.sse.connect <- client
			client.hub = hub
			log.Println("<GameHub Client Connect>")

		case client := <-hub.disconnect:
			delete(hub.clients, client)
			hub.sse.disconnect <- client
			log.Println("<GameHub Client Disonnect>")

		case message := <-hub.serverMessageChan:
			hub.handleMessage(message)
		}
	}
}

func (hub *GameHub) handleMessage(message string) {
	log.Println("<GameHub Message>:", message)
}
