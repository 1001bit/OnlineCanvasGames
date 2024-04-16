package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Layer of RT which is responsible for handling connection with GameRT SSE
type GameRTClient struct {
	gameRT *GameRT

	stopChan chan struct{}
	doneChan chan struct{}

	writer    http.ResponseWriter
	writeChan chan MessageJSON
}

func NewGameRTClient(writer http.ResponseWriter) *GameRTClient {
	return &GameRTClient{
		gameRT: nil,

		stopChan: make(chan struct{}),
		doneChan: make(chan struct{}),

		writer:    writer,
		writeChan: make(chan MessageJSON),
	}
}

// Constantly wait for message from writeChan and write it to writer
func (client *GameRTClient) Run(ctx context.Context) {
	log.Println("<GameRTClient Run>")

	defer func() {
		client.gameRT.disconnectClientChan <- client
		log.Println("<GameRTClient Run End>")
	}()

	for {
		select {
		case message := <-client.writeChan:
			// Write message to writer if server told to do so
			client.writeMessage(message)
			log.Println("<GameRTClient Write Message>")

		case <-client.doneChan:
			// When parent closed done chan
			return

		case <-client.stopChan:
			// When server asked to stop client
			return

		case <-ctx.Done():
			// When http request is done
			return
		}
	}
}

func (client *GameRTClient) Stop() {
	client.stopChan <- struct{}{}
}

func (client *GameRTClient) writeMessage(message MessageJSON) {
	messageByte, err := json.Marshal(message)
	if err != nil {
		log.Println("error marshaling GameRTClient message:", err)
		return
	}

	fmt.Fprintf(client.writer, "data: %s\n\n", messageByte)
	client.writer.(http.Flusher).Flush()
}
