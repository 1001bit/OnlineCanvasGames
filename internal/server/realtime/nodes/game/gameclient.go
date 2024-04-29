package gamenode

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/runflow"
)

// Layer of RT which is responsible for handling connection with GameNode SSE
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
	log.Println("<GameClient Run>")

	defer func() {
		client.Flow.CloseDone()
		log.Println("<GameClient Done>")
	}()

	for {
		select {
		case msg := <-client.writeChan:
			// Write message to writer if server told to do so
			client.writeMessage(msg)
			log.Println("<GameClient Write Message>")

		case <-client.Flow.Stopped():
			// When server asked to stop client
			return

		case <-ctx.Done():
			// When http request is done
			return
		}
	}
}

func (client *GameClient) writeMessage(msg *message.JSON) {
	msgByte, err := json.Marshal(msg)
	if err != nil {
		log.Println("error marshaling GameClient message:", err)
		return
	}

	fmt.Fprintf(client.writer, "data: %s\n\n", msgByte)
	client.writer.(http.Flusher).Flush()
}
