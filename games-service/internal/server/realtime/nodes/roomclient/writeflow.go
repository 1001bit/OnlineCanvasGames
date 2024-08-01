package roomclient

import (
	"log"
	"time"

	"github.com/1001bit/ocg-games-service/internal/server/message"
	"github.com/gorilla/websocket"
)

// flow that is responsible for writing messages to conn
func (client *RoomClient) writeFlow() {
	log.Println("--<RoomClient writeFlow>")
	defer log.Println("--<RoomClient writeFlow Done>")

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Ping every tick
			client.pingConn()

		case msg := <-client.writeChan:
			// Write message to conn if server told to do so
			client.writeMessageToConn(msg)

		case <-client.Flow.Done():
			// when client is done
			return
		}
	}
}

// ping connection every tick of ticker
func (client *RoomClient) pingConn() {
	client.conn.SetWriteDeadline(time.Now().Add(writeWait)) // if WriteMessage can't ping in writeWait period, client is disconnected

	// Ping the connection with special message
	err := client.conn.WriteMessage(websocket.PingMessage, nil)
	// if couldn't write message - disconnect
	if err != nil {
		log.Println("stop on ping")
		go client.Flow.Stop()
	}
}

// write message to connection
func (client *RoomClient) writeMessageToConn(msg *message.JSON) {
	client.conn.SetWriteDeadline(time.Now().Add(writeWait)) // if WriteMessage can't send message in writeWait period, client is disconnected

	err := client.conn.WriteJSON(msg)
	// if couldn't write message - disconnect
	if err != nil || msg.Type == CloseMsgType {
		go client.Flow.Stop()
	}
}
