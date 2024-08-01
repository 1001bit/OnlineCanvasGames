package gamenode

import (
	"log"

	"github.com/1001bit/onlinecanvasgames/services/games/internal/server/message"
	"github.com/1001bit/onlinecanvasgames/services/games/internal/server/realtime/nodes/gameclient"
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

			log.Println("<GameNode +Client>:", gameNode.Clients.ChildrenSet.Length())

		case client := <-gameNode.Clients.ToDisconnect():
			// When server asked to disconnect a client
			gameNode.disconnectClient(client)
			log.Println("<GameNode -Client>:", gameNode.Clients.ChildrenSet.Length())

		case <-gameNode.Flow.Done():
			return
		}
	}
}

// connect GameNode client to GameNode
func (gameNode *GameNode) connectClient(client *gameclient.GameClient) {
	gameNode.Clients.ChildrenSet.Insert(client)
}

// disconnect GameNode client from gameNode
func (gameNode *GameNode) disconnectClient(client *gameclient.GameClient) {
	gameNode.Clients.ChildrenSet.Delete(client)
}
