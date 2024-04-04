package sse

import "log"

type GameHub struct {
	SSELayer
	sse *GamesSSE
}

func NewGameHub(sse *GamesSSE) *GameHub {
	return &GameHub{
		SSELayer: MakeSSELayer(),
		sse:      sse,
	}
}

func (hub *GameHub) Run() {
	for {
		select {
		case client := <-hub.connect:
			hub.clients[client] = true
			hub.sse.connect <- client
			log.Println("<GameHub Connect>")

		case client := <-hub.disconnect:
			delete(hub.clients, client)
			hub.sse.disconnect <- client
			log.Println("<GameHub Disonnect>")

		case message := <-hub.messageChan:
			hub.handleMessage(message)
		}
	}
}

func (hub *GameHub) handleMessage(message string) {
	log.Println("<GameHub Message>:", message)
}
