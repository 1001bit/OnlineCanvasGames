package gamenode

import (
	"log"

	"github.com/neinBit/ocg-games-service/internal/server/message"
	"github.com/neinBit/ocg-games-service/internal/server/realtime/nodes/roomnode"
)

func (gameNode *GameNode) roomsFlow() {
	log.Println("--<GameNode roomsFlow>")
	defer log.Println("--<GameNode roomsFlow Done>")

	for {
		select {
		case room := <-gameNode.Rooms.ToConnect():
			// When server asked to connect a room
			gameNode.connectRoom(room)
			log.Println("<GameNode +Room>:", gameNode.Rooms.IDMap.Length())

		case room := <-gameNode.Rooms.ToDisconnect():
			// When server asked to disconnect a client
			gameNode.disconnectRoom(room)

			// update roomsJSON on room delete
			gameNode.updateRoomsJSON()

			log.Println("<GameNode -Room>:", gameNode.Rooms.IDMap.Length())

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
	gameNode.Rooms.IDMap.Set(room.GetID(), room)

	room.ConfirmConnectToGame()
}

// disconnect RoomNode from GameNode
func (gameNode *GameNode) disconnectRoom(room *roomnode.RoomNode) {
	gameNode.Rooms.IDMap.Delete(room.GetID())
}

// update gameNode.roomsJSON rooms list to send to all the clients of gameNode
func (gameNode *GameNode) updateRoomsJSON() {
	gameNode.roomsJSON = make([]RoomJSON, gameNode.Rooms.IDMap.Length())

	i := 0

	idMap, rUnlockFunc := gameNode.Rooms.IDMap.GetMapForRead()
	defer rUnlockFunc()

	for _, roomNode := range idMap {
		select {
		case <-roomNode.ConnectedToGame():
			// if room connected
		case <-roomNode.Flow.Done():
			// if room is already done
			continue
		}

		gameNode.roomsJSON[i] = RoomJSON{
			Owner:   roomNode.GetOwnerName(),
			Clients: roomNode.Clients.IDMap.Length(),
			Limit:   roomNode.GetPlayersLimit(),
			ID:      roomNode.GetID(),
		}

		i += 1
	}

	go gameNode.GlobalWriteMessage(&message.JSON{
		Type: "rooms",
		Body: gameNode.roomsJSON,
	})
}
