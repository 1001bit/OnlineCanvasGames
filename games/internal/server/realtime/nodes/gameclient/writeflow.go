package gameclient

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/neinBit/ocg-games-service/internal/server/message"
)

func (client *GameClient) writeFlow() {
	log.Println("--<GameClient writeFlow>")
	defer log.Println("--<GameClient writeFlow Done>")

	for {
		select {
		case msg := <-client.writeChan:
			// Write message to writer if server told to do so
			client.writeMessageToWriter(msg)
			log.Println("<GameClient Write Message>")

		case <-client.Flow.Done():
			return
		}
	}
}

func (client *GameClient) writeMessageToWriter(msg *message.JSON) {
	msgByte, err := json.Marshal(msg)
	if err != nil {
		log.Println("error marshaling GameClient message:", err)
		return
	}

	fmt.Fprintf(client.writer, "data: %s\n\n", msgByte)
	client.writer.(http.Flusher).Flush()
}
