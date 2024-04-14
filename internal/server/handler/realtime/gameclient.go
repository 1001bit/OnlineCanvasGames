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

	done chan struct{}

	writer    http.ResponseWriter
	writeChan chan []byte
}

func NewGameRTClient(writer http.ResponseWriter) *GameRTClient {
	return &GameRTClient{
		gameRT: nil,

		done: make(chan struct{}),

		writer:    writer,
		writeChan: make(chan []byte),
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
			log.Println("<GameRTClient Write>:", string(message))

		case <-client.done:
			return

		case <-ctx.Done():
			return
		}
	}
}

func (client *GameRTClient) writeMessage(message []byte) {
	fmt.Fprintf(client.writer, "data: %s\n\n", message)
	client.writer.(http.Flusher).Flush()
}
