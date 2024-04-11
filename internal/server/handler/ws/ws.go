package ws

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	ErrNoRooms = errors.New("no room exists")
)

type GamesWS struct {
	rooms map[int]*GameRoom

	createRoomChan   chan *GameRoom
	removeRoomIDChan chan int
}

func NewGamesWS() *GamesWS {
	ws := &GamesWS{
		rooms: make(map[int]*GameRoom),

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

	// TODO: Incorrect roomID error handling
	roomID, err := strconv.Atoi(r.PathValue("roomid"))
	if err != nil {
		return
	}

	room, ok := ws.rooms[roomID]
	if !ok {
		return
	}

	client := NewClient(conn)
	room.connectChan <- client

	go client.readPump()
	go client.writePump()
}

func (ws *GamesWS) Run() {
	log.Println("<GameWS Run>")

	for {
		select {
		case room := <-ws.createRoomChan:
			ws.createRoom(room)
			log.Println("<GameWS Room Create>")

		case roomID := <-ws.removeRoomIDChan:
			ws.removeRoomByID(roomID)
			log.Println("<GameWS Room Remove>")
		}
	}
}

func (ws *GamesWS) createRoom(room *GameRoom) {
	ws.rooms[room.id] = room
	room.ws = ws

	go room.Run()
}

func (ws *GamesWS) removeRoomByID(id int) {
	delete(ws.rooms, id)
}

func (ws *GamesWS) pickRandomRoomID() (int, error) {
	if len(ws.rooms) == 0 {
		return 0, ErrNoRooms
	}

	k := rand.Intn(len(ws.rooms))
	for roomID := range ws.rooms {
		if k == 0 {
			return roomID, nil
		}
		k--
	}

	return 0, ErrNoRooms
}
