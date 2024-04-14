package realtime

import (
	"encoding/json"
	"log"
	"strconv"
	"time"
)

type SSEMessageJSON struct {
	Type string `json:"type"`
	Data any
}

// Room that will be sent to client
type RoomJSON struct {
	Owner   string `json:"owner"`
	Clients int    `json:"clients"`
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

	globalWriteChan chan []byte

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

		globalWriteChan: make(chan []byte),

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
		case client := <-gameRT.connectClientChan:
			gameRT.connectClient(client)
			log.Println("<GameRT +Client>:", len(gameRT.clients))

		case client := <-gameRT.disconnectClientChan:
			gameRT.disconnectClient(client)
			log.Println("<GameRT -Client>:", len(gameRT.clients))

		case room := <-gameRT.connectRoomChan:
			gameRT.connectRoom(room)
			log.Println("<GameRT +Room>:", len(gameRT.rooms))
		case room := <-gameRT.disconnectRoomChan:
			gameRT.disconnectRoom(room)
			log.Println("<GameRT -Room>:", len(gameRT.rooms))

		case message := <-gameRT.globalWriteChan:
			gameRT.globalWriteMessage(message)
			log.Printf("<GameRT Global Message>: %s\n", message)

		case <-gameRT.done:
			return
		}
	}
}

// connect GameRT client to GameRT
func (gameRT *GameRT) connectClient(client *GameRTClient) {
	gameRT.clients[client] = true
	client.gameRT = gameRT

	str, err := json.Marshal(gameRT.prepareRoomsMessage())
	if err != nil {
		log.Println("error marshaling rooms:", err)
	}
	client.writeChan <- str
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

	// notify all the clients about new room
	str, err := json.Marshal(gameRT.prepareRoomsMessage())
	if err != nil {
		log.Println("error marshaling rooms:", err)
	}
	gameRT.globalWriteChan <- str
}

// disconnect RoomRT from GameRT
func (gameRT *GameRT) disconnectRoom(room *RoomRT) {
	if _, ok := gameRT.rooms[room.id]; !ok {
		return
	}

	delete(gameRT.rooms, room.id)
	close(room.done)

	// notify all the clients about room delete
	str, err := json.Marshal(gameRT.prepareRoomsMessage())
	if err != nil {
		log.Println("error marshaling rooms:", err)
	}
	gameRT.globalWriteChan <- str
}

// write a message to every client
func (gameRT *GameRT) globalWriteMessage(message []byte) {
	for client := range gameRT.clients {
		client.writeChan <- message
	}
}

// get all the rooms in JSON
func (gameRT *GameRT) prepareRoomsMessage() SSEMessageJSON {
	message := SSEMessageJSON{
		Type: "rooms",
	}

	rooms := make([]RoomJSON, len(gameRT.rooms))

	for i := 0; i < len(gameRT.rooms); i++ {
		roomRT := gameRT.rooms[i]
		rooms[i] = RoomJSON{
			// TODO: Make username instead of userid
			Owner:   strconv.Itoa(roomRT.owner.userID),
			Clients: len(roomRT.clients),
		}
	}

	message.Data = rooms

	return message
}
