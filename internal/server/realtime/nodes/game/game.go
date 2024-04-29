package gamenode

import (
	"errors"
	"log"

	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/children"
	roomnode "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/room"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/runflow"
)

var ErrNoRooms = errors.New("no rooms in the game")

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
	Clients children.Children[GameClient]

	roomsJSON           *message.JSON
	roomsJSONUpdateChan chan struct{}

	globalWriteChan chan *message.JSON

	gameID int
}

func NewGameNode(id int) *GameNode {
	return &GameNode{
		Flow: runflow.MakeRunFlow(),

		Rooms:   children.MakeChildrenWithID[roomnode.RoomNode](),
		Clients: children.MakeChildren[GameClient](),

		roomsJSON: &message.JSON{
			Type: "rooms",
			Body: make([]RoomJSON, 0),
		},
		roomsJSONUpdateChan: make(chan struct{}),

		globalWriteChan: make(chan *message.JSON),

		gameID: id,
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
			go func() {
				client.writeChan <- gameNode.roomsJSON
			}()

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

		case msg := <-gameNode.globalWriteChan:
			// Write message to every client if server told to do so
			gameNode.globalWriteMessage(msg)
			log.Println("<GameNode Global Message>")

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

// ask gameNode to update gameNode.roomsJSON
func (gameNode *GameNode) RequestUpdatingRoomsJSON() {
	gameNode.roomsJSONUpdateChan <- struct{}{}
}

func (gameNode *GameNode) GetID() int {
	return gameNode.gameID
}

// connect GameNode client to GameNode
func (gameNode *GameNode) connectClient(client *GameClient) {
	gameNode.Clients.ChildMap[client] = true
}

// disconnect GameNode client from gameNode
func (gameNode *GameNode) disconnectClient(client *GameClient) {
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

// write a message to every client
func (gameNode *GameNode) globalWriteMessage(msg *message.JSON) {
	for client := range gameNode.Clients.ChildMap {
		client.writeChan <- msg
	}
}

// update gameNode.roomsJSON rooms list to send to all the clients of gameNode
func (gameNode *GameNode) updateRoomsJSON() {
	roomsJSON := make([]RoomJSON, 0)

	for _, roomNode := range gameNode.Rooms.IDMap {
		<-roomNode.ConnectedToGame()

		roomsJSON = append(roomsJSON, RoomJSON{
			Owner:   roomNode.GetOwnerName(),
			Clients: len(roomNode.Clients.ChildMap),
			ID:      roomNode.GetID(),
		})
	}

	gameNode.roomsJSON.Body = roomsJSON
	gameNode.globalWriteMessage(gameNode.roomsJSON)
}
