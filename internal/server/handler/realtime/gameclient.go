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

	done chan struct{}

	writer    http.ResponseWriter
	writeChan chan MessageJSON
}

func NewGameRTClient(writer http.ResponseWriter) *GameRTClient {
	return &GameRTClient{
		gameRT: nil,

		done: make(chan struct{}),

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
		// Write message to writer if server told to do so
		case message := <-client.writeChan:
			client.writeMessage(message)
			log.Println("<GameRTClient Write Message>")

		// When game closed client.done
		case <-client.done:
			return

		// When http request is done
		case <-ctx.Done():
			return
		}
	}
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
