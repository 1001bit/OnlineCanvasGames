package basenode

import (
	"context"
	"net/http"

	gamenode "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/game"
)

// handle SSE endpoint
func (baseRT *BaseRT) ConnectToGame(ctx context.Context, w http.ResponseWriter, gameID int) error {
	gameRT, ok := baseRT.games.IDMap[gameID]
	if !ok {
		return ErrNoGame
	}

	client := gamenode.NewGameRTClient(w)

	// RUN GameRTClient
	go func() {
		gameRT.Clients.ConnectChild(client)
		client.Run(ctx)
		gameRT.Clients.DisconnectChild(client)
	}()

	<-client.Flow.Done()

	return nil
}
