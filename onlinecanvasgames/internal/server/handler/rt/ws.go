package rt

import (
	"net/http"
	"strconv"

	"github.com/1001bit/OnlineCanvasGames/internal/auth/claimscontext"
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/basenode"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/roomclient"
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
	// GameID
	gameID, err := strconv.Atoi(r.PathValue("gameid"))
	if err != nil {
		closeConnWithMessage(conn, "Wrong game id!")
		return
	}

	// RoomID
	roomID, err := strconv.Atoi(r.PathValue("roomid"))
	if err != nil {
		closeConnWithMessage(conn, "Wrong room id!")
		return
	}

	// UserID and username
	userID, username, err := claimscontext.GetClaims(r.Context())
	if err != nil {
		closeConnWithMessage(conn, "Unauthorized")
		return
	}

	err = baseNode.ConnectToRoom(conn, gameID, roomID, userID, username)
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
