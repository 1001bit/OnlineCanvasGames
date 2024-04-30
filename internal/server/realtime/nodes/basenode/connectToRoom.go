package basenode

import (
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/roomclient"
	"github.com/gorilla/websocket"
)

// Handle WS endpoint
func (baseNode *BaseNode) ConnectToRoom(conn *websocket.Conn, gameID, roomID, userID int, userName string) error {
	gameNode, ok := baseNode.games.IDMap[gameID]
	if !ok {
		return ErrNoGame
	}

	roomNode, ok := gameNode.Rooms.IDMap[roomID]
	if !ok {
		return ErrNoRoom
	}

	user := roomclient.RoomClientUser{
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
