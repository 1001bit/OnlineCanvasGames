package gamenode

import (
	"log"

	gamemodel "github.com/1001bit/OnlineCanvasGames/internal/model/game"
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/children"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/gameclient"
	roomnode "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/roomnode"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/runflow"
)

// Room that will be sent to client
type RoomJSON struct {
	Owner   string `json:"owner"`
	Clients int    `json:"clients"`
	ID      int    `json:"id"`
}

// Layer of RT which is responsible for game hub and containing rooms
type GameNode struct {
	Flow runflow.RunFlow

	Rooms   children.ChildrenWithID[roomnode.RoomNode]
	Clients children.Children[gameclient.GameClient]

	roomsJSON           []RoomJSON
	roomsJSONUpdateChan chan struct{}

	game gamemodel.Game
}

func NewGameNode(game gamemodel.Game) *GameNode {
	return &GameNode{
		Flow: runflow.MakeRunFlow(),

		Rooms:   children.MakeChildrenWithID[roomnode.RoomNode](),
		Clients: children.MakeChildren[gameclient.GameClient](),

		roomsJSON:           make([]RoomJSON, 0),
		roomsJSONUpdateChan: make(chan struct{}),

		game: game,
	}
}

func (gameNode *GameNode) Run() {
	log.Println("<GameNode Run>")

	defer func() {
		gameNode.Flow.CloseDone()

		log.Println("<GameNode Done>")
	}()

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

		case <-gameNode.Flow.Stopped():
			// When server asked to stop running
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

	gameNode.GlobalWriteMessage(&message.JSON{
		Type: "rooms",
		Body: gameNode.roomsJSON,
	})
}
