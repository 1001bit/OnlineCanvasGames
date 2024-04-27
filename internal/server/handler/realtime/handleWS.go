package realtime

import (
	"net/http"
	"strconv"

	"github.com/1001bit/OnlineCanvasGames/internal/auth"
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
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

func closeConnWithMessage(conn *websocket.Conn, text string) {
	conn.WriteJSON(message.JSON{
		Type: "message",
		Body: text,
	})
}

// Handle WS endpoint
func (rt *Realtime) HandleRoomWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get game from path
	gameID, err := strconv.Atoi(r.PathValue("gameid"))
	if err != nil {
		closeConnWithMessage(conn, "Wrong game id!")
		return
	}
	game, ok := rt.games.idMap[gameID]
	if !ok {
		closeConnWithMessage(conn, "Wrong game id!")
		return
	}

	// Get room from path
	roomID, err := strconv.Atoi(r.PathValue("roomid"))
	if err != nil {
		closeConnWithMessage(conn, "Wrong room id!")
		return
	}

	room, ok := game.rooms.idMap[roomID]
	if !ok {
		closeConnWithMessage(conn, "Wrong room id!")
		return
	}

	// Get user from JWT
	claims, err := auth.JWTClaimsByRequest(r)
	if err != nil {
		closeConnWithMessage(conn, "Unauthorized!")
		return
	}
	user := RoomClientUser{}

	// ID
	userIDfloat, ok := claims["userID"].(float64) // for some reason, in JWT it's stored as float64
	if !ok {
		closeConnWithMessage(conn, "Unauthorized!")
		return
	}
	user.id = int(userIDfloat)

	// Name
	user.name, ok = claims["username"].(string)
	if !ok {
		closeConnWithMessage(conn, "Unauthorized!")
		return
	}

	// Create client and start client
	client := NewRoomClient(conn, user)
	go client.Run(room)
}
