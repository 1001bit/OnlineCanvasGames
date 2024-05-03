package basenode

import (
	rtclient "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/client"
	rterror "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/error"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/roomclient"
	"github.com/gorilla/websocket"
)

// Handle WS endpoint
func (baseNode *BaseNode) ConnectToRoom(conn *websocket.Conn, gameID, roomID, userID int, userName string) error {
	gameNode, ok := baseNode.games.IDMap[gameID]
	if !ok {
		return rterror.ErrNoGame
	}

	roomNode, ok := gameNode.Rooms.IDMap[roomID]
	if !ok {
		return rterror.ErrNoRoom
	}

	user := rtclient.User{
		ID:   userID,
		Name: userName,
	}

	// Create client and start client
	client := roomclient.NewRoomClient(conn, user)

	// RUN RoomClient
	go func() {
		roomNode.Clients.ConnectChild(client)
		baseNode.roomsClients.ConnectChild(client)

		client.Run(roomNode)

		roomNode.Clients.DisconnectChild(client)
		baseNode.roomsClients.DisconnectChild(client)
	}()

	return nil
}
