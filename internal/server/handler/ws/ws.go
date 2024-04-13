package ws

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/1001bit/OnlineCanvasGames/internal/auth"
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
	rooms              map[int]*GameRoom
	connectRoomChan    chan *GameRoom
	disconnectRoomChan chan *GameRoom

	clients              map[int]*Client
	connectClientChan    chan *Client
	disconnectClientChan chan *Client
}

func NewGamesWS() *GamesWS {
	ws := &GamesWS{
		rooms:              make(map[int]*GameRoom),
		connectRoomChan:    make(chan *GameRoom),
		disconnectRoomChan: make(chan *GameRoom),

		clients:              make(map[int]*Client),
		connectClientChan:    make(chan *Client),
		disconnectClientChan: make(chan *Client),
	}

	return ws
}

// upgrade connection from http to ws and connect client to requested room
func (ws *GamesWS) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	// TODO: server error handling
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

	claims, err := auth.JWTClaimsByRequest(r)
	// TODO: Incorrect token error handling
	if err != nil {
		return
	}
	userIDstr, ok := claims["userID"]
	if !ok {
		return
	}
	userID := int(userIDstr.(float64)) // for some reason, it's stored in float64
	client := NewClient(conn, userID)

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
			log.Println("<GameWS +Room>:", len(ws.rooms))

		case room := <-ws.disconnectRoomChan:
			ws.disconnectRoom(room)
			log.Println("<GameWS -Room>:", len(ws.rooms))

		case client := <-ws.connectClientChan:
			ws.connectClient(client)
			log.Println("<GameWS +Client>:", len(ws.clients))

		case client := <-ws.disconnectClientChan:
			ws.disconnectClient(client)
			log.Println("<GameWS -Client>:", len(ws.clients))
		}
	}
}

// TODO: Put context here
// Create new room, connect it to ws, return it
func (ws *GamesWS) ConnectNewRoom() *GameRoom {
	newRoom := NewGameRoom()

	ws.connectRoomChan <- newRoom
	// wait until room is fully connected and initialized
	newRoom.waitUntilConnectedToWS()

	return newRoom
}

// connect a room to ws
func (ws *GamesWS) connectRoom(room *GameRoom) {
	room.id = int(time.Now().UnixMicro())
	ws.rooms[room.id] = room
	room.ws = ws

	room.confirmConnectToWS()

	go room.Run()
}

// disconnect a room from ws by it's id
func (ws *GamesWS) disconnectRoom(room *GameRoom) {
	delete(ws.rooms, room.id)
}

// if client with such id is already connected, disconnect them. Add new client to list
func (ws *GamesWS) connectClient(client *Client) {
	if oldClient, ok := ws.clients[client.userID]; ok {
		oldClient.room.disconnectClientChan <- oldClient
	}

	ws.clients[client.userID] = client
}

// if client with requested id is requested to disconnect and the same client exists - delete
func (ws *GamesWS) disconnectClient(client *Client) {
	if oldClient, ok := ws.clients[client.userID]; ok && oldClient == client {
		delete(ws.clients, client.userID)
	}
}

// TODO: Make it thread safe
// pick a random room and return it's id
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
