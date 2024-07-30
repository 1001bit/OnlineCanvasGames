package basenode

import (
	"context"
	"net/http"

	"github.com/neinBit/ocg-games-service/internal/server/realtime/nodes/gameclient"
)

// handle SSE endpoint
func (baseNode *BaseNode) ConnectToGame(ctx context.Context, w http.ResponseWriter, gameID int) error {
	gameNode, ok := baseNode.games.IDMap.Get(gameID)
	if !ok {
		return ErrNoGame
	}

	client := gameclient.NewGameClient(w)

	// RUN GameClient
	go func() {
		gameNode.Clients.ConnectChild(client, gameNode.Flow.Done())
		client.Run(ctx)
		gameNode.Clients.DisconnectChild(client, gameNode.Flow.Done())
	}()

	select {
	case <-client.Flow.Done():
	case <-ctx.Done():
	}

	return nil
}