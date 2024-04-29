package basenode

import (
	"context"
	"net/http"

	gamenode "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/game"
)

// handle SSE endpoint
func (baseNode *BaseNode) ConnectToGame(ctx context.Context, w http.ResponseWriter, gameID int) error {
	gameNode, ok := baseNode.games.IDMap[gameID]
	if !ok {
		return ErrNoGame
	}

	client := gamenode.NewGameClient(w)

	// RUN GameClient
	go func() {
		gameNode.Clients.ConnectChild(client)
		client.Run(ctx)
		gameNode.Clients.DisconnectChild(client)
	}()

	<-client.Flow.Done()

	return nil
}
