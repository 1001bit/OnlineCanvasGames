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

type RoomClientUser struct {
	ID   int
	Name string
}

// Layer of RT which is responsible for handling connection WS
type RoomClient struct {
	roomRT *RoomRT

	stopChan chan struct{}
	doneChan chan struct{}

	conn *websocket.Conn
	user RoomClientUser

	writeChan chan MessageJSON
	readChan  chan MessageJSON
}

func NewRoomClient(conn *websocket.Conn) *RoomClient {
	return &RoomClient{
		roomRT: nil,

		stopChan: make(chan struct{}),
		doneChan: make(chan struct{}),

		conn: conn,
		user: RoomClientUser{
			ID:   0,
			Name: "",
		},

		writeChan: make(chan MessageJSON),
		readChan:  make(chan MessageJSON),
	}
}

func (client *RoomClient) Run() {
	log.Println("<RoomClient Run>")

	// ticker that indicates the need to send ping message
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		if client.roomRT != nil {
			client.roomRT.disconnectClientChan <- client
		}
		ticker.Stop()
		client.conn.Close()
		log.Println("<RoomClient Run End>")
	}()

	go client.readPump()

	for {
		select {
		case <-ticker.C:
			// Ping every tick
			client.pingConn()

		case message := <-client.writeChan:
			// Write message to conn if server told to do so
			client.writeMessage(message)

		case message := <-client.readChan:
			// Handle messages that were read in readPump
			client.handleReadMessage(message)

		case <-client.stopChan:
			// when server asked to stop running
			return

		case <-client.doneChan:
			// when parent closed doneChan
			return
		}
	}
}

func (client *RoomClient) Stop() {
	client.stopChan <- struct{}{}
}

// constantly read messages from connection
func (client *RoomClient) readPump() {
	log.Println("<RoomClient ReadPump>")

	defer func() {
		client.Stop()
		log.Println("<RoomClient ReadPump End>")
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
			case <-client.doneChan:
				return
			}
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
		client.stopWithMessage("Unexpected error!")
	}
}

// write message to connection
func (client *RoomClient) writeMessage(message MessageJSON) {
	client.conn.SetWriteDeadline(time.Now().Add(writeWait)) // if WriteMessage can't send message in writeWait period, client is disconnected

	messageByte, err := json.Marshal(message)
	if err != nil {
		log.Println("err marshaling RoomClient message:", err)
		return
	}

	err = client.conn.WriteMessage(websocket.TextMessage, messageByte)
	// if couldn't write message - disconnect
	if err != nil {
		client.stopWithMessage("Unexpected error!")
	}
}

// process read message
func (client *RoomClient) handleReadMessage(message MessageJSON) {
	// simply tell room about read message
	client.roomRT.readChan <- RoomReadMessage{
		client:  client,
		message: message,
	}
}

// send message to client and close after
func (client *RoomClient) stopWithMessage(text string) {
	client.writeChan <- MessageJSON{
		Type: "message",
		Body: text,
	}
	client.Stop()
}
