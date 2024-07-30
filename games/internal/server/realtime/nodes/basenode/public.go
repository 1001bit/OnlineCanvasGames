package basenode

import (
	"context"
	"time"

	"github.com/neinBit/ocg-games-service/internal/gamemodel"
	"github.com/neinBit/ocg-games-service/internal/server/realtime/nodes/gamenode"
	"github.com/neinBit/ocg-games-service/internal/server/realtime/nodes/roomnode"
)

// create new room and connect it to BaseNode
func (baseNode *BaseNode) ConnectNewRoom(ctx context.Context, gameID int) (*roomnode.RoomNode, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	gameNode, ok := baseNode.games.IDMap.Get(gameID)
	if !ok {
		return nil, ErrNoGame
	}

	room := roomnode.NewRoomNode(gameID)

	// RUN roomNode
	go func() {
		gameNode.Rooms.ConnectChild(room, gameNode.Flow.Done())
		room.Run(gameNode)
		gameNode.Rooms.DisconnectChild(room, gameNode.Flow.Done())
	}()

	// wait until room connected to BaseNode
	select {
	case <-room.ConnectedToGame():
		return room, nil

	case <-baseNode.Flow.Done():
		go room.Flow.Stop()
		return nil, ErrCreateRoom

	case <-ctx.Done():
		go room.Flow.Stop()
		return nil, ErrCreateRoom
	}
}

func (baseNode *BaseNode) GetGameByID(id int) (*gamenode.GameNode, bool) {
	return baseNode.games.IDMap.Get(id)
}

func (baseNode *BaseNode) GetGamesJSON() []gamemodel.Game {
	return baseNode.gamesJSON
}
