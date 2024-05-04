package gamenode

import (
	"log"

	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/roomnode"
)

func (gameNode *GameNode) roomsFlow() {
	log.Println("--<GameNode roomsFlow>")
	defer log.Println("--<GameNode roomsFlow Done>")

	for {
		select {
		case room := <-gameNode.Rooms.ToConnect():
			// When server asked to connect a room
			gameNode.connectRoom(room)
			log.Println("<GameNode +Room>:", len(gameNode.Rooms.IDMap))

		case room := <-gameNode.Rooms.ToDisconnect():
			// When server asked to disconnect a client
			gameNode.disconnectRoom(room)

			// update roomsJSON on room delete
			gameNode.updateRoomsJSON()

			log.Println("<GameNode -Room>:", len(gameNode.Rooms.IDMap))

		case <-gameNode.roomsJSONUpdateChan:
			// When server asked to update roomsJSON
			gameNode.updateRoomsJSON()
			log.Println("<GameNode RoomsJSON Update>")

		case <-gameNode.Flow.Done():
			return
		}
	}
}

// connect RoomNode to GameNode
func (gameNode *GameNode) connectRoom(room *roomnode.RoomNode) {
	room.SetRandomID()
	gameNode.Rooms.IDMap[room.GetID()] = room

	room.ConfirmConnectToGame()
}

// disconnect RoomNode from GameNode
func (gameNode *GameNode) disconnectRoom(room *roomnode.RoomNode) {
	delete(gameNode.Rooms.IDMap, room.GetID())
}

// update gameNode.roomsJSON rooms list to send to all the clients of gameNode
func (gameNode *GameNode) updateRoomsJSON() {
	gameNode.roomsJSON = make([]RoomJSON, 0)
	for _, roomNode := range gameNode.Rooms.IDMap {
		<-roomNode.ConnectedToGame()

		gameNode.roomsJSON = append(gameNode.roomsJSON, RoomJSON{
			Owner:   roomNode.GetOwnerName(),
			Clients: len(roomNode.Clients.IDMap),
			ID:      roomNode.GetID(),
		})
	}

	go gameNode.GlobalWriteMessage(&message.JSON{
		Type: "rooms",
		Body: gameNode.roomsJSON,
	})
}
