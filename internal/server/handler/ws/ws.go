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
	rooms                map[int]*GameRoom
	connectRoomChan      chan *GameRoom
	disconnectRoomIDChan chan int
}

func NewGamesWS() *GamesWS {
	ws := &GamesWS{
		rooms: make(map[int]*GameRoom),

		connectRoomChan:      make(chan *GameRoom),
		disconnectRoomIDChan: make(chan int),
	}

	return ws
}

func (ws *GamesWS) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error upgrading connection:", err)
		return
	}

	// TODO: Restrict access for users that are already connected
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
	room.connectClientChan <- client

	go client.readPump()
	go client.writePump()
}

func (ws *GamesWS) Run() {
	log.Println("<GameWS Run>")
	defer log.Println("<GameWS Run End>")

	for {
		select {
		case room := <-ws.connectRoomChan:
			ws.connectRoom(room)
			log.Println("<GameWS Room Connect>")

		case roomID := <-ws.disconnectRoomIDChan:
			ws.disconnectRoomByID(roomID)
			log.Println("<GameWS Room Disconnect>")
		}
	}
}

// TODO: Put context here
func (ws *GamesWS) ConnectNewRoom() *GameRoom {
	newRoom := NewGameRoom()
	ws.connectRoomChan <- newRoom
	return newRoom
}

func (ws *GamesWS) connectRoom(room *GameRoom) {
	ws.rooms[room.id] = room
	room.ws = ws

	go room.Run()
}

func (ws *GamesWS) disconnectRoomByID(id int) {
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
