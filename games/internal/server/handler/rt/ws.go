package rt

import (
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/neinBit/ocg-games-service/internal/server/message"
	"github.com/neinBit/ocg-games-service/internal/server/realtime/nodes/basenode"
	"github.com/neinBit/ocg-games-service/internal/server/realtime/nodes/roomclient"
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
		Type: roomclient.CloseMsgType,
		Body: text,
	})
}

func HandleRoomWS(w http.ResponseWriter, r *http.Request, baseNode *basenode.BaseNode) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get params from path
	// Game title
	gameTitle := r.PathValue("title")

	// RoomID
	roomID, err := strconv.Atoi(r.PathValue("roomid"))
	if err != nil {
		closeConnWithMessage(conn, "Wrong room id!")
		return
	}

	// username
	username := r.Header.Get("X-Username")

	if username == "" {
		closeConnWithMessage(conn, "Unauthorized!")
		return
	}

	err = baseNode.ConnectToRoom(conn, gameTitle, roomID, username)
	switch err {
	case nil:
		// no error
	case basenode.ErrNoGame:
		closeConnWithMessage(conn, "Wrong game id!")
	case basenode.ErrNoRoom:
		closeConnWithMessage(conn, "Wrong room id!")
	default:
		closeConnWithMessage(conn, "Unexpected error!")
	}
}
