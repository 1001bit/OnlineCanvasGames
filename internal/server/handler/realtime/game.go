package realtime

import (
	"log"
	"strconv"
	"time"
)

// Room that will be sent to client
type RoomJSON struct {
	Owner   string `json:"owner"`
	Clients int    `json:"clients"`
	ID      int    `json:"id"`
}

// Layer of RT which is responsible for game hub and containing rooms
type GameRT struct {
	rt *Realtime

	done chan struct{}

	rooms              map[int]*RoomRT
	connectRoomChan    chan *RoomRT
	disconnectRoomChan chan *RoomRT

	clients              map[*GameRTClient]bool
	connectClientChan    chan *GameRTClient
	disconnectClientChan chan *GameRTClient

	globalWriteChan chan MessageJSON

	gameID int
}

func NewGameRT() *GameRT {
	return &GameRT{
		rt: nil,

		done: make(chan struct{}),

		rooms:              make(map[int]*RoomRT),
		connectRoomChan:    make(chan *RoomRT),
		disconnectRoomChan: make(chan *RoomRT),

		clients:              make(map[*GameRTClient]bool),
		connectClientChan:    make(chan *GameRTClient),
		disconnectClientChan: make(chan *GameRTClient),

		globalWriteChan: make(chan MessageJSON),

		gameID: 0,
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
		// Clients
		case client := <-gameRT.connectClientChan:
			gameRT.connectClient(client)

			// send rooms data to client on it's join
			client.writeChan <- gameRT.prepareRoomsMessage()

			log.Println("<GameRT +Client>:", len(gameRT.clients))

		case client := <-gameRT.disconnectClientChan:
			gameRT.disconnectClient(client)
			log.Println("<GameRT -Client>:", len(gameRT.clients))

		// Rooms
		case room := <-gameRT.connectRoomChan:
			gameRT.connectRoom(room)
			log.Println("<GameRT +Room>:", len(gameRT.rooms))

		case room := <-gameRT.disconnectRoomChan:
			gameRT.disconnectRoom(room)
			log.Println("<GameRT -Room>:", len(gameRT.rooms))

		// Global messages
		case message := <-gameRT.globalWriteChan:
			gameRT.globalWriteMessage(message)
			log.Println("<GameRT Global Message>")

		// Done
		case <-gameRT.done:
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
	close(client.done)
}

// connect RoomRT to GameRT
func (gameRT *GameRT) connectRoom(room *RoomRT) {
	room.id = int(time.Now().UnixMicro())
	gameRT.rooms[room.id] = room
	room.gameRT = gameRT

	close(room.connectedToRT)

	go room.Run()
}

// disconnect RoomRT from GameRT
func (gameRT *GameRT) disconnectRoom(room *RoomRT) {
	if _, ok := gameRT.rooms[room.id]; !ok {
		return
	}

	delete(gameRT.rooms, room.id)
	close(room.done)
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
		ownerUserID := 0
		if roomRT.owner != nil {
			ownerUserID = roomRT.owner.userID
		}

		rooms = append(rooms, RoomJSON{
			// TODO: Make username instead of userid
			Owner:   strconv.Itoa(ownerUserID),
			Clients: len(roomRT.clients),
			ID:      roomRT.id,
		})
	}

	message.Body = rooms

	return message
}
