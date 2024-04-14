package realtime

import (
	"log"
	"time"
)

// Layer of RT which is responsible for game hub and containing rooms
type GameRT struct {
	rt *Realtime

	rooms              map[int]*RoomRT
	connectRoomChan    chan *RoomRT
	disconnectRoomChan chan *RoomRT

	clients              map[*GameRTClient]bool
	connectClientChan    chan *GameRTClient
	disconnectClientChan chan *GameRTClient

	gameID int
}

func NewGameRT() *GameRT {
	return &GameRT{
		rt: nil,

		rooms:              make(map[int]*RoomRT),
		connectRoomChan:    make(chan *RoomRT),
		disconnectRoomChan: make(chan *RoomRT),

		clients:              make(map[*GameRTClient]bool),
		connectClientChan:    make(chan *GameRTClient),
		disconnectClientChan: make(chan *GameRTClient),

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
	close(client.writeChan)
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
	delete(gameRT.rooms, room.id)
}
