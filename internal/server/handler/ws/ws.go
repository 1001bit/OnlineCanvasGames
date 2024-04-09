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
		WSLayer: MakeWSLayer(),

		rooms:   make(map[int]*GameRoom),
		clients: make(map[*Client]bool),

		createRoomChan:   make(chan *GameRoom),
		removeRoomIDChan: make(chan int),
	}

	return ws
}

func (ws *GamesWS) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error upgrading connection:", err)
		return
	}

	client := NewClient(conn)

	// TODO: Remove this when request-based room connection will be made
	randomRoom := ws.rooms[ws.pickRandomRoomID()]
	randomRoom.connectChan <- client

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
			log.Println("<GameWS Client Connect>")

		case client := <-ws.disconnectChan:
			delete(ws.clients, client)
			log.Println("<GameWS Client Disconnect>")

		// Room
		case room := <-ws.createRoomChan:
			// TODO: Replace to room creation interface
			room.id = int(time.Now().Unix())

			room.ws = ws
			ws.rooms[room.id] = room
			go room.Run()
			log.Println("<GameWS Room Create>")

		case roomID := <-ws.removeRoomIDChan:
			delete(ws.rooms, roomID)
			log.Println("<GameWS Room Remove>")

		// Message from client
		case message := <-ws.clientMessageChan:
			ws.handleClientMessage(message)
		}
	}
}

func (ws *GamesWS) handleClientMessage(message ClientMessage) {
	log.Println("<GameWS Message>:", message)
}

func (ws *GamesWS) pickRandomRoomID() int {
	if len(ws.clients) == 0 {
		return 0
	}

	k := rand.Intn(len(ws.clients))
	for roomID := range ws.rooms {
		if k == 0 {
			return roomID
		}
		k--
	}
	return 0
}
