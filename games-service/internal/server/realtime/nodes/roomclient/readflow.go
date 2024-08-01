package roomclient

import (
	"io"
	"log"
	"time"

	"github.com/1001bit/ocg-games-service/internal/server/message"
	rtclient "github.com/1001bit/ocg-games-service/internal/server/realtime/client"
	"github.com/gorilla/websocket"
)

// flow that is responsible for reading messages from conn
func (client *RoomClient) readFlow(r RoomNodeReader) {
	log.Println("--<RoomClient readFlow>")
	defer log.Println("--<RoomClient readFlow Done>")

	// On Pong
	client.conn.SetReadDeadline(time.Now().Add(pongWait)) // if ReadMessage doesn't get any message after pongWait period, readPump stops
	client.conn.SetPongHandler(func(string) error {       // updates pongWait period
		client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		// Get msg from client
		msg := &message.JSON{}
		err := client.conn.ReadJSON(msg)

		switch err {
		case nil:
			// no error

		case io.ErrUnexpectedEOF:
			// unmarshaling error, probably bad message from client
			continue

		default:
			// reading error
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("error reading ws message:", err)
			}

			client.Flow.Stop()
			return
		}

		select {
		case <-client.Flow.Done():
			return
		default:
			// continue read loop
		}

		client.handleReadMessage(msg, r)
	}
}

// process read message
func (client *RoomClient) handleReadMessage(msg *message.JSON, roomNodeReader RoomNodeReader) {
	// simply tell room about read message
	go roomNodeReader.ReadMessage(rtclient.MessageWithClient{
		Client:  client,
		Message: msg,
	})
}
