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

// Layer of RT which is responsible for handling connection with GameRT SSE
type GameRTClient struct {
	Flow runflow.RunFlow

	writer    http.ResponseWriter
	writeChan chan *message.JSON
}

func NewGameRTClient(writer http.ResponseWriter) *GameRTClient {
	return &GameRTClient{
		Flow: runflow.MakeRunFlow(),

		writer:    writer,
		writeChan: make(chan *message.JSON),
	}
}

// Constantly wait for message from writeChan and write it to writer
func (client *GameRTClient) Run(ctx context.Context) {
	log.Println("<GameRTClient Run>")

	defer func() {
		client.Flow.CloseDone()
		log.Println("<GameRTClient Done>")
	}()

	for {
		select {
		case msg := <-client.writeChan:
			// Write message to writer if server told to do so
			client.writeMessage(msg)
			log.Println("<GameRTClient Write Message>")

		case <-client.Flow.Stopped():
			// When server asked to stop client
			return

		case <-ctx.Done():
			// When http request is done
			return
		}
	}
}

func (client *GameRTClient) writeMessage(msg *message.JSON) {
	msgByte, err := json.Marshal(msg)
	if err != nil {
		log.Println("error marshaling GameRTClient message:", err)
		return
	}

	fmt.Fprintf(client.writer, "data: %s\n\n", msgByte)
	client.writer.(http.Flusher).Flush()
}
