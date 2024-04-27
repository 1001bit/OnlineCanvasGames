package gamenode

import (
	"errors"
	"log"

	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/realtime/children"
	roomnode "github.com/1001bit/OnlineCanvasGames/internal/server/handler/realtime/nodes/room"
	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/realtime/runflow"
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
)

var ErrNoRooms = errors.New("no rooms in the game")

// Room that will be sent to client
type RoomJSON struct {
	Owner   string `json:"owner"`
	Clients int    `json:"clients"`
	ID      int    `json:"id"`
}

// Layer of RT which is responsible for game hub and containing rooms
type GameRT struct {
	Flow runflow.RunFlow

	Rooms   children.ChildrenWithID[roomnode.RoomRT]
	Clients children.Children[GameRTClient]

	roomsJSON           *message.JSON
	roomsJSONUpdateChan chan struct{}

	globalWriteChan chan *message.JSON

	gameID int
}

func NewGameRT(id int) *GameRT {
	return &GameRT{
		Flow: runflow.MakeRunFlow(),

		Rooms:   children.MakeChildrenWithID[roomnode.RoomRT](),
		Clients: children.MakeChildren[GameRTClient](),

		roomsJSON: &message.JSON{
			Type: "rooms",
			Body: make([]RoomJSON, 0),
		},
		roomsJSONUpdateChan: make(chan struct{}),

		globalWriteChan: make(chan *message.JSON),

		gameID: id,
	}
}

func (gameRT *GameRT) Run() {
	log.Println("<GameRT Run>")

	defer func() {
		gameRT.Flow.CloseDone()

		log.Println("<GameRT Done>")
	}()

	for {
		select {
		case client := <-gameRT.Clients.ToConnect():
			// When server asked to connect a client
			gameRT.connectClient(client)

			// send roomsJSON to client on it's join
			go func() {
				client.writeChan <- gameRT.roomsJSON
			}()

			log.Println("<GameRT +Client>:", len(gameRT.Clients.ChildMap))

		case client := <-gameRT.Clients.ToDisconnect():
			// When server asked to disconnect a client
			gameRT.disconnectClient(client)
			log.Println("<GameRT -Client>:", len(gameRT.Clients.ChildMap))

		case room := <-gameRT.Rooms.ToConnect():
			// When server asked to connect a room
			gameRT.connectRoom(room)
			log.Println("<GameRT +Room>:", len(gameRT.Rooms.IDMap))

		case room := <-gameRT.Rooms.ToDisconnect():
			// When server asked to disconnect a client
			gameRT.disconnectRoom(room)

			// update roomsJSON on room delete
			gameRT.updateRoomsJSON()

			log.Println("<GameRT -Room>:", len(gameRT.Rooms.IDMap))

		case msg := <-gameRT.globalWriteChan:
			// Write message to every client if server told to do so
			gameRT.globalWriteMessage(msg)
			log.Println("<GameRT Global Message>")

		case <-gameRT.roomsJSONUpdateChan:
			// When server asked to update roomsJSON
			gameRT.updateRoomsJSON()
			log.Println("<GameRT RoomsJSON Update>")

		case <-gameRT.Flow.Stopped():
			// When server asked to stop running
			return
		}
	}
}

// ask gameRT to update gameRT.roomsJSON
func (gameRT *GameRT) RequestUpdatingRoomsJSON() {
	gameRT.roomsJSONUpdateChan <- struct{}{}
}

func (gameRT *GameRT) GetID() int {
	return gameRT.gameID
}

// connect GameRT client to GameRT
func (gameRT *GameRT) connectClient(client *GameRTClient) {
	gameRT.Clients.ChildMap[client] = true
}

// disconnect GameRT client from gameRT
func (gameRT *GameRT) disconnectClient(client *GameRTClient) {
	delete(gameRT.Clients.ChildMap, client)
}

// connect RoomRT to GameRT
func (gameRT *GameRT) connectRoom(room *roomnode.RoomRT) {
	room.SetRandomID()
	gameRT.Rooms.IDMap[room.GetID()] = room

	room.ConfirmConnectToGame()
}

// disconnect RoomRT from GameRT
func (gameRT *GameRT) disconnectRoom(room *roomnode.RoomRT) {
	delete(gameRT.Rooms.IDMap, room.GetID())
}

// write a message to every client
func (gameRT *GameRT) globalWriteMessage(msg *message.JSON) {
	for client := range gameRT.Clients.ChildMap {
		client.writeChan <- msg
	}
}

// update gameRT.roomsJSON rooms list to send to all the clients of gameRT
func (gameRT *GameRT) updateRoomsJSON() {
	roomsJSON := make([]RoomJSON, 0)

	for _, roomRT := range gameRT.Rooms.IDMap {
		<-roomRT.ConnectedToGame()

		roomsJSON = append(roomsJSON, RoomJSON{
			Owner:   roomRT.GetOwnerName(),
			Clients: len(roomRT.Clients.ChildMap),
			ID:      roomRT.GetID(),
		})
	}

	gameRT.roomsJSON.Body = roomsJSON
	gameRT.globalWriteMessage(gameRT.roomsJSON)
}
