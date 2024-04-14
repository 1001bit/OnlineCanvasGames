package realtime

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
	pongWait  = 60 * time.Second
	// must be shorter that pong wait, because ping should be sent before, than read deadline happens
	pingPeriod = pongWait * 9 / 10
)

// Layer of RT which is responsible for handling connection WS
type RoomRTClient struct {
	roomRT *RoomRT

	done chan struct{}

	conn      *websocket.Conn
	userID    int
	writeChan chan []byte
}

func NewRoomRTClient(conn *websocket.Conn, userID int) *RoomRTClient {
	return &RoomRTClient{
		roomRT: nil,

		done: make(chan struct{}),

		conn:      conn,
		userID:    userID,
		writeChan: make(chan []byte),
	}
}

// close websocket connection
func (client *RoomRTClient) closeConn() {
	closeConnWithMessage(client.conn, "closed websocket connection")
}

// constantly read messages from connection
func (client *RoomRTClient) readPump() {
	log.Println("<RoomRTClient ReadPump>")

	defer func() {
		client.roomRT.disconnectClientChan <- client
		client.closeConn()
		log.Println("<RoomRTClient ReadPump End>")
	}()

	// On Pong
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		// Get message from client
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("error reading ws message:", err)
			}
			break
		}

		// send newly read the message to room
		client.roomRT.readChan <- RoomRTClientMessage{
			client: client,
			text:   string(message),
		}
	}
}

// constantly check messages in writeChan and send them to connection
func (client *RoomRTClient) writePump() {
	log.Println("<RoomRTClient WritePump>")

	// ticker that indicates the need to send ping message
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		client.closeConn()
		log.Println("<RoomRTClient WritePump End>")
	}()

	for {
		select {
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// Ping the connection when ticker worked
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case message := <-client.writeChan:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))

			// write message from writeChan to connection
			client.conn.WriteMessage(websocket.TextMessage, message)

			log.Printf("<RoomRTClient Write>: %s\n", message)

		case <-client.done:
			return
		}
	}
}
