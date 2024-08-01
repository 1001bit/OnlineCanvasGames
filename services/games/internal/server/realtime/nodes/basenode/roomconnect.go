package basenode

import (
	rtclient "github.com/1001bit/onlinecanvasgames/services/games/internal/server/realtime/client"
	"github.com/1001bit/onlinecanvasgames/services/games/internal/server/realtime/nodes/roomclient"
	"github.com/gorilla/websocket"
)

// Handle WS endpoint
func (baseNode *BaseNode) ConnectToRoom(conn *websocket.Conn, gameTitle string, roomID int, username string) error {
	gameNode, ok := baseNode.games.ChildrenMap.Get(gameTitle)
	if !ok {
		return ErrNoGame
	}

	roomNode, ok := gameNode.Rooms.ChildrenMap.Get(roomID)
	if !ok {
		return ErrNoRoom
	}

	user := rtclient.User{
		Name: username,
	}

	// Create client and start client
	client := roomclient.NewRoomClient(conn, user)

	// RUN RoomClient
	go func() {
		roomNode.Clients.ConnectChild(client, roomNode.Flow.Done())
		baseNode.roomsClients.ConnectChild(client, baseNode.Flow.Done())

		client.Run(roomNode)

		roomNode.Clients.DisconnectChild(client, roomNode.Flow.Done())
		baseNode.roomsClients.DisconnectChild(client, baseNode.Flow.Done())
	}()

	return nil
}
