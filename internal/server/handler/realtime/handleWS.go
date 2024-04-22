package realtime

import (
	"net/http"
	"strconv"

	"github.com/1001bit/OnlineCanvasGames/internal/auth"
	"github.com/gorilla/websocket"
)

type WSMessage struct {
	Message string `json:"message"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Handle WS endpoint
func (rt *Realtime) HandleRoomWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create client and start client
	client := NewRoomClient(conn)
	go client.Run()

	// Get game from path
	gameID, err := strconv.Atoi(r.PathValue("gameid"))
	if err != nil {
		client.stopWithMessage("Wrong game id!")
		return
	}
	game, ok := rt.games[gameID]
	if !ok {
		client.stopWithMessage("Wrong game id!")
		return
	}

	// Get room from path
	roomID, err := strconv.Atoi(r.PathValue("roomid"))
	if err != nil {
		client.stopWithMessage("Wrong room id!")
		return
	}

	room, ok := game.rooms[roomID]
	if !ok {
		client.stopWithMessage("Wrong room id!")
		return
	}
	// connect client to the room
	room.connectClientChan <- client

	// Get user from JWT
	claims, err := auth.JWTClaimsByRequest(r)
	if err != nil {
		client.stopWithMessage("Unauthorized!")
		return
	}

	// ID
	userIDstr, ok := claims["userID"].(float64) // for some reason, in JWT it's stored as float64
	if !ok {
		client.stopWithMessage("Unauthorized!")
		return
	}
	client.user.id = int(userIDstr)

	// Name
	client.user.name, ok = claims["username"].(string)
	if !ok {
		client.stopWithMessage("Unauthorized!")
		return
	}
}
