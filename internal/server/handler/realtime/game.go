package realtime

import (
	"errors"
	"log"
	"math/rand"
	"time"
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

	stopChan chan struct{}
	doneChan chan struct{}

	rooms              map[int]*RoomRT
	connectRoomChan    chan *RoomRT
	disconnectRoomChan chan *RoomRT

	clients              map[*GameRTClient]bool
	connectClientChan    chan *GameRTClient
	disconnectClientChan chan *GameRTClient

	globalWriteChan chan MessageJSON

	gameID int
}

func NewGameRT(id int) *GameRT {
	return &GameRT{
		rt: nil,

		stopChan: make(chan struct{}),
		doneChan: make(chan struct{}),

		rooms:              make(map[int]*RoomRT),
		connectRoomChan:    make(chan *RoomRT),
		disconnectRoomChan: make(chan *RoomRT),

		clients:              make(map[*GameRTClient]bool),
		connectClientChan:    make(chan *GameRTClient),
		disconnectClientChan: make(chan *GameRTClient),

		globalWriteChan: make(chan MessageJSON),

		gameID: id,
	}
}

func (gameRT *GameRT) Run() {
	log.Println("<GameRT Run>")

	defer func() {
		gameRT.rt.disconnectGameChan <- gameRT
		log.Println("<GameRT Run End>")
	}()

	for {
		select {
		case client := <-gameRT.connectClientChan:
			// When server asked to connect a client
			gameRT.connectClient(client)

			// send rooms data to client on it's join
			client.writeChan <- gameRT.prepareRoomsMessage()

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

			// send rooms json data globally on room delete
			go gameRT.globallyWriteRoomsMessage()

			log.Println("<GameRT -Room>:", len(gameRT.rooms))

		case message := <-gameRT.globalWriteChan:
			// Write message to every client if server told to do so
			gameRT.globalWriteMessage(message)
			log.Println("<GameRT Global Message>")

		case <-gameRT.stopChan:
			// When server asked to stop running
			return

		case <-gameRT.doneChan:
			// When parent closed done chan
			return
		}
	}
}

func (gameRT *GameRT) Stop() {
	gameRT.stopChan <- struct{}{}
}

// returns random room
func (gameRT *GameRT) PickRandomRoom() (*RoomRT, error) {
	if len(gameRT.rooms) == 0 {
		return nil, ErrNoRooms
	}

	k := rand.Intn(len(gameRT.rooms))
	for _, room := range gameRT.rooms {
		if k == 0 {
			return room, nil
		}
		k--
	}
	return nil, ErrNoRooms
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
	close(client.doneChan)
}

// connect RoomRT to GameRT
func (gameRT *GameRT) connectRoom(room *RoomRT) {
	room.id = int(time.Now().UnixMicro())
	gameRT.rooms[room.id] = room
	room.gameRT = gameRT

	close(room.connectedToRT)
}

// disconnect RoomRT from GameRT
func (gameRT *GameRT) disconnectRoom(room *RoomRT) {
	if _, ok := gameRT.rooms[room.id]; !ok {
		return
	}

	delete(gameRT.rooms, room.id)
	close(room.doneChan)
}

// write a message to every client
func (gameRT *GameRT) globalWriteMessage(message MessageJSON) {
	for client := range gameRT.clients {
		client.writeChan <- message
	}
}

// get all the rooms in JSON format
func (gameRT *GameRT) prepareRoomsMessage() MessageJSON {
	message := MessageJSON{
		Type: "rooms",
	}

	rooms := make([]RoomJSON, 0)

	for _, roomRT := range gameRT.rooms {
		<-roomRT.connectedToRT
		ownerUser := RoomClientUser{}
		if roomRT.owner != nil {
			ownerUser = roomRT.owner.user
		}

		rooms = append(rooms, RoomJSON{
			Owner:   ownerUser.Name,
			Clients: len(roomRT.clients),
			ID:      roomRT.id,
		})
	}

	message.Body = rooms

	return message
}

func (gameRT *GameRT) globallyWriteRoomsMessage() {
	gameRT.globalWriteChan <- gameRT.prepareRoomsMessage()
}
