package realtime

import (
	"net/http"
	"strconv"

	"github.com/1001bit/OnlineCanvasGames/internal/auth"
	"github.com/gorilla/websocket"
)

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

	// Get game from path
	gameID, err := strconv.Atoi(r.PathValue("gameid"))
	if err != nil {
		closeConnWithMessage(conn, "wrong game id!")
		return
	}
	game, ok := rt.games[gameID]
	if !ok {
		closeConnWithMessage(conn, "wrong game id!")
		return
	}

	// Get room from path
	roomID, err := strconv.Atoi(r.PathValue("roomid"))
	if err != nil {
		closeConnWithMessage(conn, "wrong room id!")
		return
	}

	room, ok := game.rooms[roomID]
	if !ok {
		closeConnWithMessage(conn, "wrong room id!")
		return
	}

	// Get userID from JWT
	claims, err := auth.JWTClaimsByRequest(r)
	if err != nil {
		closeConnWithMessage(conn, "unauthorized!")
		return
	}
	userIDstr, ok := claims["userID"]
	if !ok {
		closeConnWithMessage(conn, "unauthorized!")
		return
	}

	userID := int(userIDstr.(float64)) // for some reason, it's stored in float64

	// Connect client to the room and start client
	client := NewRoomRTClient(conn, userID)

	room.connectClientChan <- client

	go client.readPump()
	go client.writePump()
}

func closeConnWithMessage(conn *websocket.Conn, text string) {
	conn.WriteMessage(websocket.CloseMessage, []byte("  "+text))
	conn.Close()
}
