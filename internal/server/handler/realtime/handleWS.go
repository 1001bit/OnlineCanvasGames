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
	client := NewRoomRTClient(conn)

	// Get game from path
	gameID, err := strconv.Atoi(r.PathValue("gameid"))
	if err != nil {
		client.closeConnWithMessage("Wrong game id!")
		return
	}
	game, ok := rt.games[gameID]
	if !ok {
		client.closeConnWithMessage("Wrong game id!")
		return
	}

	// Get room from path
	roomID, err := strconv.Atoi(r.PathValue("roomid"))
	if err != nil {
		client.closeConnWithMessage("Wrong room id!")
		return
	}

	room, ok := game.rooms[roomID]
	if !ok {
		client.closeConnWithMessage("Wrong room id!")
		return
	}
	// connect client to the room
	room.connectClientChan <- client

	// Get userID from JWT
	claims, err := auth.JWTClaimsByRequest(r)
	if err != nil {
		client.closeConnWithMessage("Unauthorized!")
		return
	}
	userIDstr, ok := claims["userID"]
	if !ok {
		client.closeConnWithMessage("Unauthorized!")
		return
	}

	// set client userID from JWT
	client.userID = int(userIDstr.(float64)) // for some reason, it's stored in float64
}
