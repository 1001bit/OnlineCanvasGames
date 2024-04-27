package basenode

import (
	"net/http"
	"strconv"

	"github.com/1001bit/OnlineCanvasGames/internal/auth"
	roomnode "github.com/1001bit/OnlineCanvasGames/internal/server/handler/realtime/nodes/room"
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
	"github.com/gorilla/websocket"
)

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
func (baseRT *BaseRT) HandleRoomWS(w http.ResponseWriter, r *http.Request) {
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
	gameRT, ok := baseRT.games.IDMap[gameID]
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

	roomRT, ok := gameRT.Rooms.IDMap[roomID]
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
	user := roomnode.RoomClientUser{}

	// ID
	userIDfloat, ok := claims["userID"].(float64) // for some reason, in JWT it's stored as float64
	if !ok {
		closeConnWithMessage(conn, "Unauthorized!")
		return
	}
	user.ID = int(userIDfloat)

	// Name
	user.Name, ok = claims["username"].(string)
	if !ok {
		closeConnWithMessage(conn, "Unauthorized!")
		return
	}

	// Create client and start client
	client := roomnode.NewRoomClient(conn, user)

	// RUN RoomClient
	go func() {
		roomRT.Clients.ConnectChild(client)
		client.Run(roomRT)
		roomRT.Clients.DisconnectChild(client)
	}()
}
