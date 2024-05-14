package basenode

import (
	"context"
	"net/http"

	rterror "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/error"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/gameclient"
)

// handle SSE endpoint
func (baseNode *BaseNode) ConnectToGame(ctx context.Context, w http.ResponseWriter, gameID int) error {
	gameNode, ok := baseNode.games.IDMap[gameID]
	if !ok {
		return rterror.ErrNoGame
	}

	client := gameclient.NewGameClient(w)

	// RUN GameClient
	go func() {
		gameNode.Clients.ConnectChild(client)
		client.Run(ctx)
		gameNode.Clients.DisconnectChild(client)
	}()

	select {
	case <-client.Flow.Done():
	case <-ctx.Done():
	}

	return nil
}