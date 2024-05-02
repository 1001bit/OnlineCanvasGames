package gamenode

import (
	"log"

	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/gameclient"
)

func (gameNode *GameNode) clientsFlow() {
	log.Println("--<GameNode clientsFlow>")
	defer log.Println("--<GameNode clientsFlow Done>")

	for {
		select {
		case client := <-gameNode.Clients.ToConnect():
			// When server asked to connect a client
			gameNode.connectClient(client)

			// send roomsJSON to client on it's join
			go client.WriteMessage(&message.JSON{
				Type: "rooms",
				Body: gameNode.roomsJSON,
			})

			log.Println("<GameNode +Client>:", len(gameNode.Clients.ChildMap))

		case client := <-gameNode.Clients.ToDisconnect():
			// When server asked to disconnect a client
			gameNode.disconnectClient(client)
			log.Println("<GameNode -Client>:", len(gameNode.Clients.ChildMap))

		case <-gameNode.Flow.Done():
			return
		}
	}
}

// connect GameNode client to GameNode
func (gameNode *GameNode) connectClient(client *gameclient.GameClient) {
	gameNode.Clients.ChildMap[client] = true
}

// disconnect GameNode client from gameNode
func (gameNode *GameNode) disconnectClient(client *gameclient.GameClient) {
	delete(gameNode.Clients.ChildMap, client)
}
