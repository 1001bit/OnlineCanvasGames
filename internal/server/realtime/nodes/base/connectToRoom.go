package basenode

import (
	roomnode "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/room"
	"github.com/gorilla/websocket"
)

// Handle WS endpoint
func (baseRT *BaseRT) ConnectToRoom(conn *websocket.Conn, gameID, roomID, userID int, userName string) error {
	gameRT, ok := baseRT.games.IDMap[gameID]
	if !ok {
		return ErrNoGame
	}

	roomRT, ok := gameRT.Rooms.IDMap[roomID]
	if !ok {
		return ErrNoRoom
	}

	user := roomnode.RoomClientUser{
		ID:   userID,
		Name: userName,
	}

	// Create client and start client
	client := roomnode.NewRoomClient(conn, user)

	// RUN RoomClient
	go func() {
		roomRT.Clients.ConnectChild(client)
		baseRT.roomsClients.ConnectChild(client)

		client.Run(roomRT)

		roomRT.Clients.DisconnectChild(client)
		baseRT.roomsClients.DisconnectChild(client)
	}()

	return nil
}
