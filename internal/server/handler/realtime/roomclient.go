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

	conn   *websocket.Conn
	userID int

	writeChan chan MessageJSON
	readChan  chan MessageJSON
}

func NewRoomRTClient(conn *websocket.Conn) *RoomRTClient {
	return &RoomRTClient{
		roomRT: nil,

		done: make(chan struct{}),

		conn:   conn,
		userID: 0,

		writeChan: make(chan MessageJSON),
		readChan:  make(chan MessageJSON),
	}
}

func (client *RoomRTClient) Run() {
	log.Println("<RoomRTClient Run>")

	// ticker that indicates the need to send ping message
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		client.conn.Close()
		log.Println("<RoomRTClient Run End>")
	}()

	go client.readPump()

	for {
		select {
		// Ping every tick
		case <-ticker.C:
			client.pingConn()

		// Write message to conn if server told to do so
		case message := <-client.writeChan:
			client.writeMessage(message)

		// Handle messages that were read in readPump
		case message := <-client.readChan:
			client.handleReadMessage(message)

		// When room closed client.done chan
		case <-client.done:
			// warn client about ws closure
			client.writeMessage(MessageJSON{
				Type: "message",
				Body: "WebSocket connection close!",
			})
			return
		}
	}
}

// constantly read messages from connection
func (client *RoomRTClient) readPump() {
	log.Println("<RoomRTClient ReadPump>")

	defer func() {
		if client.roomRT != nil {
			client.roomRT.disconnectClientChan <- client
		}
		log.Println("<RoomRTClient ReadPump End>")
	}()

	// On Pong
	client.conn.SetReadDeadline(time.Now().Add(pongWait)) // if ReadMessage doesn't get any message after pongWait period, readPump stops
	client.conn.SetPongHandler(func(string) error {       // updates pongWait period
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

		// transform message into struct and throw into channel
		messageStruct := MessageJSON{}
		err = json.Unmarshal(message, &messageStruct)
		if err == nil {
			select {
			case client.readChan <- messageStruct:
				// send message to read chan
			case <-client.done:
				return
			}
		}
	}
}

// ping connection every tick of ticker
func (client *RoomRTClient) pingConn() {
	client.conn.SetWriteDeadline(time.Now().Add(writeWait)) // if WriteMessage can't ping in writeWait period, client is disconnected

	// Ping the connection with special message
	err := client.conn.WriteMessage(websocket.PingMessage, nil)
	// if couldn't write message - disconnect
	if err != nil {
		client.roomRT.disconnectClientChan <- client
	}
}

// write message to connection
func (client *RoomRTClient) writeMessage(message MessageJSON) {
	client.conn.SetWriteDeadline(time.Now().Add(writeWait)) // if WriteMessage can't send message in writeWait period, client is disconnected

	messageByte, err := json.Marshal(message)
	if err != nil {
		log.Println("err marshaling RoomRTClient message:", err)
		return
	}

	err = client.conn.WriteMessage(websocket.TextMessage, messageByte)
	// if couldn't write message - disconnect
	if err != nil {
		client.roomRT.disconnectClientChan <- client
	}
}

// process read message
func (client *RoomRTClient) handleReadMessage(message MessageJSON) {
	// simply tell room about read message
	client.roomRT.readChan <- RoomReadMessage{
		client:  client,
		message: message,
	}
}

// send message to client and close after
func (client *RoomRTClient) closeConnWithMessage(text string) {
	client.writeChan <- MessageJSON{
		Type: "message",
		Body: text,
	}

	// TODO: Make something without using sleep. For example, closeMessage chan
	// wait for message to be sent, then close connection
	time.Sleep(time.Second)

	client.conn.Close()
}
