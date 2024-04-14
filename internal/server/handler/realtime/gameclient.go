package realtime

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// Layer of RT which is responsible for handling connection with GameRT SSE
type GameRTClient struct {
	gameRT *GameRT

	writer    http.ResponseWriter
	writeChan chan string
}

func NewGameRTClient(writer http.ResponseWriter) *GameRTClient {
	return &GameRTClient{
		gameRT: nil,

		writer:    writer,
		writeChan: make(chan string),
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
		case message, ok := <-client.writeChan:
			// if hub closed client.write chan
			if !ok {
				return
			}

			client.writeMessage(message)
			log.Println("<GameRTClient Write>:", string(message))

		case <-ctx.Done():
			return
		}
	}
}

func (client *GameRTClient) writeMessage(message string) {
	fmt.Fprintf(client.writer, "data: %s\n\n", message)
	client.writer.(http.Flusher).Flush()
}
