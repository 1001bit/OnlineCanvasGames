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
	rt *Realtime

	flow RunFlow

	rooms              map[int]*RoomRT
	connectRoomChan    chan *RoomRT
	disconnectRoomChan chan *RoomRT

	clients              map[*GameRTClient]bool
	connectClientChan    chan *GameRTClient
	disconnectClientChan chan *GameRTClient

	roomsJSON           *message.JSON
	roomsJSONUpdateChan chan struct{}

	globalWriteChan chan *message.JSON

	gameID int
}

func NewGameRT(id int) *GameRT {
	return &GameRT{
		rt: nil,

		flow: MakeRunFlow(),

		rooms:              make(map[int]*RoomRT),
		connectRoomChan:    make(chan *RoomRT),
		disconnectRoomChan: make(chan *RoomRT),

		clients:              make(map[*GameRTClient]bool),
		connectClientChan:    make(chan *GameRTClient),
		disconnectClientChan: make(chan *GameRTClient),

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
		gameRT.rt.disconnectGameChan <- gameRT
		gameRT.flow.CloseDone()

		log.Println("<GameRT Done>")
	}()

	for {
		select {
		case client := <-gameRT.connectClientChan:
			// When server asked to connect a client
			gameRT.connectClient(client)

			// send roomsJSON to client on it's join
			go func() {
				client.writeChan <- gameRT.roomsJSON
			}()

			log.Println("<GameRT +Client>:", len(gameRT.clients))

		case client := <-gameRT.disconnectClientChan:
			// When server asked to disconnect a client
			gameRT.disconnectClient(client)
			log.Println("<GameRT -Client>:", len(gameRT.clients))

		case room := <-gameRT.connectRoomChan:
			// When server asked to connect a room
			gameRT.connectRoom(room)
			log.Println("<GameRT +Room>:", len(gameRT.rooms))

		case room := <-gameRT.disconnectRoomChan:
			// When server asked to disconnect a client
			gameRT.disconnectRoom(room)

			// update roomsJSON on room delete
			gameRT.updateRoomsJSON()

			log.Println("<GameRT -Room>:", len(gameRT.rooms))

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
	gameRT.clients[client] = true
	client.gameRT = gameRT
}

// disconnect GameRT client from gameRT
func (gameRT *GameRT) disconnectClient(client *GameRTClient) {
	if _, ok := gameRT.clients[client]; !ok {
		return
	}

	delete(gameRT.clients, client)
}

// connect RoomRT to GameRT
func (gameRT *GameRT) connectRoom(room *RoomRT) {
	room.id = int(time.Now().UnixMicro())
	gameRT.rooms[room.id] = room
	room.gameRT = gameRT

	close(room.connectedToGame)
}

// disconnect RoomRT from GameRT
func (gameRT *GameRT) disconnectRoom(room *RoomRT) {
	if _, ok := gameRT.rooms[room.id]; !ok {
		return
	}

	delete(gameRT.rooms, room.id)
}

// write a message to every client
func (gameRT *GameRT) globalWriteMessage(msg *message.JSON) {
	for client := range gameRT.clients {
		client.writeChan <- msg
	}
}

// update gameRT.roomsJSON rooms list to send to all the clients of gameRT
func (gameRT *GameRT) updateRoomsJSON() {
	roomsJSON := make([]RoomJSON, 0)

	for _, roomRT := range gameRT.rooms {
		<-roomRT.connectedToGame

		roomOwnerName := "nobody"
		if roomRT.owner != nil {
			roomOwnerName = roomRT.owner.user.name
		}

		roomsJSON = append(roomsJSON, RoomJSON{
			Owner:   roomOwnerName,
			Clients: len(roomRT.clients),
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
