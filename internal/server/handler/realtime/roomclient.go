package realtime

import (
	"encoding/json"
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
	writeChan chan MessageJSON
}

func NewRoomRTClient(conn *websocket.Conn) *RoomRTClient {
	return &RoomRTClient{
		roomRT: nil,

		done: make(chan struct{}),

		conn:      conn,
		userID:    0,
		writeChan: make(chan MessageJSON),
	}
}

// constantly read messages from connection
func (client *RoomRTClient) readPump() {
	log.Println("<RoomRTClient ReadPump>")

	defer func() {
		if client.roomRT != nil {
			client.roomRT.disconnectClientChan <- client
		}
		client.conn.Close()
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

		// read message and send it to room
		go client.sendMsgToRoom(message)
	}
}

// transform byte message into struct and pass it to room readChan
func (client *RoomRTClient) sendMsgToRoom(message []byte) {
	messageStruct := MessageJSON{}
	err := json.Unmarshal(message, &messageStruct)
	if err != nil {
		return
	}

	// send newly read the message to room
	client.roomRT.readChan <- RoomReadMessage{
		client:  client,
		message: messageStruct,
	}
}

// constantly check messages in writeChan and send them to connection
func (client *RoomRTClient) writePump() {
	log.Println("<RoomRTClient WritePump>")

	// ticker that indicates the need to send ping message
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		client.conn.Close()
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
			client.writeMsgToConn(message)

			log.Println("<RoomRTClient Write Message>")

		case <-client.done:
			client.writeMsgToConn(MessageJSON{
				Type: "message",
				Body: "WebSocket connection close!",
			})
			return
		}
	}
}

// convert struct message to string message and send it to user through connection
func (client *RoomRTClient) writeMsgToConn(message MessageJSON) {
	messageStr, err := json.Marshal(message)
	if err != nil {
		log.Println("err marshaling RoomRTClient message:", err)
		return
	}
	log.Println(message)

	client.conn.WriteMessage(websocket.TextMessage, messageStr)
}

func (client *RoomRTClient) closeConnWithMessage(text string) {
	client.writeChan <- MessageJSON{
		Type: "message",
		Body: text,
	}

	// TODO: Make something with it. For example, closeMessage chan
	// wait for message to be sent, then close connection
	time.Sleep(time.Second)

	client.conn.Close()
}
