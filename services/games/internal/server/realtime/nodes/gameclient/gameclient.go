package gameclient

import (
	"context"
	"log"
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/games/internal/server/realtime/runflow"
	"github.com/1001bit/onlinecanvasgames/services/games/pkg/message"
)

// Layer of RT which is responsible for handling connections: GameClient > User, GameNode > GameClient
type GameClient struct {
	Flow runflow.RunFlow

	writer    http.ResponseWriter
	writeChan chan *message.JSON
}

func NewGameClient(writer http.ResponseWriter) *GameClient {
	return &GameClient{
		Flow: runflow.MakeRunFlow(),

		writer:    writer,
		writeChan: make(chan *message.JSON),
	}
}

// Constantly wait for message from writeChan and write it to writer
func (client *GameClient) Run(ctx context.Context) {
	defer client.Flow.CloseDone()

	log.Println("-<GameClient Run>")
	defer log.Println("-<GameClient Run Done>")

	go client.writeFlow()

	select {
	case <-client.Flow.Stopped():
		// When server asked to stop client
		return

	case <-ctx.Done():
		// When http request is done
		return
	}
}
