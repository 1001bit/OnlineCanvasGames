package basenode

import (
	"context"
	"time"

	gamemodel "github.com/1001bit/OnlineCanvasGames/internal/model/game"
	rterror "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/error"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/gamenode"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/roomnode"
)

// create new room and connect it to BaseNode
func (baseNode *BaseNode) ConnectNewRoom(ctx context.Context, gameID int) (*roomnode.RoomNode, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	gameNode, ok := baseNode.games.IDMap[gameID]
	if !ok {
		return nil, rterror.ErrNoGame
	}

	room := roomnode.NewRoomNode(gameID)

	// RUN roomNode
	go func() {
		gameNode.Rooms.ConnectChild(room)
		room.Run(gameNode)
		gameNode.Rooms.DisconnectChild(room)
	}()

	// wait until room connected to BaseNode
	select {
	case <-room.ConnectedToGame():
		return room, nil

	case <-baseNode.Flow.Done():
		go room.Flow.Stop()
		return nil, rterror.ErrCreateRoom

	case <-ctx.Done():
		go room.Flow.Stop()
		return nil, rterror.ErrCreateRoom
	}
}

func (baseNode *BaseNode) GetGameByID(id int) (*gamenode.GameNode, bool) {
	game, ok := baseNode.games.IDMap[id]
	return game, ok
}

func (baseNode *BaseNode) GetGamesJSON() []gamemodel.Game {
	return baseNode.gamesJSON
}
