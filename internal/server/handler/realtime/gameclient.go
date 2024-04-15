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
func (client *GameRTClient) writePump(ctx context.Context) {
	log.Println("<GameRTClient WritePump>")

	defer func() {
		client.gameRT.disconnectClientChan <- client
		log.Println("<GameRTClient WritePump End>")
	}()

	for {
		select {
		case message := <-client.writeChan:
			client.writeMessage(message)
			log.Println("<GameRTClient Write Message>")

		case <-client.done:
			return

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
