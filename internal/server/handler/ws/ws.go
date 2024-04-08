package ws

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type GamesWS struct {
	WSLayer
	rooms   map[int]*GameRoom
	clients map[*Client]bool

	createRoomChan   chan *GameRoom
	removeRoomIDChan chan int
}

func NewGamesWS() *GamesWS {
	ws := &GamesWS{
		rooms:   make(map[int]*GameRoom),
		clients: make(map[*Client]bool),

		WSLayer: MakeWSLayer(),
	}
	ws.rooms[0] = &GameRoom{}

	return ws
}

func (ws *GamesWS) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error upgrading connection:", err)
		return
	}

	room := ws.rooms[0]
	client := NewClient(conn, room)
	client.room.connectChan <- client

	go client.readPump()
	go client.writePump()
}

func (ws *GamesWS) Run() {
	log.Println("<GameWS Run>")

	for {
		select {
		// Client
		case client := <-ws.connectChan:
			ws.clients[client] = true
			log.Println("<GameWS Connect>")

		case client := <-ws.disconnectChan:
			delete(ws.clients, client)
			log.Println("<GameWS Disconnect>")

		// Room
		case room := <-ws.createRoomChan:
			room.id = int(time.Now().Unix())
			ws.rooms[room.id] = room
			log.Println("<GameWS Create Room>")

		case roomID := <-ws.removeRoomIDChan:
			delete(ws.rooms, roomID)
			log.Println("<GameWS Create Room>")

		// Message
		case message := <-ws.messageChan:
			ws.handleMessage(message)
		}
	}
}

func (ws *GamesWS) handleMessage(message string) {
	log.Println("<GameWS Message>:", message)
}

func (ws *GamesWS) PickRandomRoomID() int {
	if len(ws.clients) == 0 {
		return -1
	}

	k := rand.Intn(len(ws.clients))
	for roomID := range ws.rooms {
		if k == 0 {
			return roomID
		}
		k--
	}
	return -1
}
