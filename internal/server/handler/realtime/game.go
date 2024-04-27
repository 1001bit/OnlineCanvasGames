package realtime

import (
	"errors"
	"log"
	"time"

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
	flow RunFlow

	rooms   ChildrenWithID[RoomRT]
	clients Children[GameRTClient]

	roomsJSON           *message.JSON
	roomsJSONUpdateChan chan struct{}

	globalWriteChan chan *message.JSON

	gameID int
}

func NewGameRT(id int) *GameRT {
	return &GameRT{
		flow: MakeRunFlow(),

		rooms:   MakeChildrenWithID[RoomRT](),
		clients: MakeChildren[GameRTClient](),

		roomsJSON: &message.JSON{
			Type: "rooms",
			Body: make([]RoomJSON, 0),
		},
		roomsJSONUpdateChan: make(chan struct{}),

		globalWriteChan: make(chan *message.JSON),

		gameID: id,
	}
}

func (gameRT *GameRT) Run(rt *Realtime) {
	log.Println("<GameRT Run>")

	rt.games.ConnectChild(gameRT)

	defer func() {
		go rt.games.DisconnectChild(gameRT)
		gameRT.flow.CloseDone()

		log.Println("<GameRT Done>")
	}()

	for {
		select {
		case client := <-gameRT.clients.connectChan:
			// When server asked to connect a client
			gameRT.connectClient(client)

			// send roomsJSON to client on it's join
			go func() {
				client.writeChan <- gameRT.roomsJSON
			}()

			log.Println("<GameRT +Client>:", len(gameRT.clients.childMap))

		case client := <-gameRT.clients.disconnectChan:
			// When server asked to disconnect a client
			gameRT.disconnectClient(client)
			log.Println("<GameRT -Client>:", len(gameRT.clients.childMap))

		case room := <-gameRT.rooms.connectChan:
			// When server asked to connect a room
			gameRT.connectRoom(room)
			log.Println("<GameRT +Room>:", len(gameRT.rooms.idMap))

		case room := <-gameRT.rooms.disconnectChan:
			// When server asked to disconnect a client
			gameRT.disconnectRoom(room)

			// update roomsJSON on room delete
			gameRT.updateRoomsJSON()

			log.Println("<GameRT -Room>:", len(gameRT.rooms.idMap))

		case msg := <-gameRT.globalWriteChan:
			// Write message to every client if server told to do so
			gameRT.globalWriteMessage(msg)
			log.Println("<GameRT Global Message>")

		case <-gameRT.roomsJSONUpdateChan:
			// When server asked to update roomsJSON
			gameRT.updateRoomsJSON()
			log.Println("<GameRT RoomsJSON Update>")

		case <-gameRT.flow.Stopped():
			// When server asked to stop running
			return
		}
	}
}

// connect GameRT client to GameRT
func (gameRT *GameRT) connectClient(client *GameRTClient) {
	gameRT.clients.childMap[client] = true
}

// disconnect GameRT client from gameRT
func (gameRT *GameRT) disconnectClient(client *GameRTClient) {
	delete(gameRT.clients.childMap, client)
}

// connect RoomRT to GameRT
func (gameRT *GameRT) connectRoom(room *RoomRT) {
	room.id = int(time.Now().UnixMicro())
	gameRT.rooms.idMap[room.id] = room

	close(room.connectedToGame)
}

// disconnect RoomRT from GameRT
func (gameRT *GameRT) disconnectRoom(room *RoomRT) {
	delete(gameRT.rooms.idMap, room.id)
}

// write a message to every client
func (gameRT *GameRT) globalWriteMessage(msg *message.JSON) {
	for client := range gameRT.clients.childMap {
		client.writeChan <- msg
	}
}

// update gameRT.roomsJSON rooms list to send to all the clients of gameRT
func (gameRT *GameRT) updateRoomsJSON() {
	roomsJSON := make([]RoomJSON, 0)

	for _, roomRT := range gameRT.rooms.idMap {
		<-roomRT.connectedToGame

		roomOwnerName := "nobody"
		if roomRT.owner != nil {
			roomOwnerName = roomRT.owner.user.name
		}

		roomsJSON = append(roomsJSON, RoomJSON{
			Owner:   roomOwnerName,
			Clients: len(roomRT.clients.childMap),
			ID:      roomRT.id,
		})
	}

	gameRT.roomsJSON.Body = roomsJSON
	gameRT.globalWriteMessage(gameRT.roomsJSON)
}

// ask gameRT to update gameRT.roomsJSON
func (gameRT *GameRT) requestUpdatingRoomsJSON() {
	gameRT.roomsJSONUpdateChan <- struct{}{}
}
