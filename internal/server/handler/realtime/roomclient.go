package realtime

import (
	"io"
	"log"
	"time"

	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
	pongWait  = 60 * time.Second
	// must be shorter that pong wait, because ping should be sent before, than read deadline happens
	pingPeriod = pongWait * 9 / 10
)

type RoomClientUser struct {
	id   int
	name string
}

// Layer of RT which is responsible for handling connection WS
type RoomClient struct {
	flow RunFlow

	conn *websocket.Conn
	user RoomClientUser

	writeChan chan *message.JSON
	readChan  chan *message.JSON
}

func NewRoomClient(conn *websocket.Conn, user RoomClientUser) *RoomClient {
	return &RoomClient{
		flow: MakeRunFlow(),

		conn: conn,
		user: user,

		writeChan: make(chan *message.JSON),
		readChan:  make(chan *message.JSON),
	}
}

func (client *RoomClient) Run(roomRT *RoomRT) {
	log.Println("<RoomClient Run>")

	// ticker that indicates the need to send ping message
	ticker := time.NewTicker(pingPeriod)

	roomRT.clients.ConnectChild(client)

	defer func() {
		go roomRT.clients.DisconnectChild(client)

		ticker.Stop()
		client.conn.Close()

		client.flow.CloseDone()

		log.Println("<RoomClient Done>")
	}()

	go client.readPump()

	for {
		select {
		case <-ticker.C:
			// Ping every tick
			client.pingConn()

		case msg := <-client.writeChan:
			// Write message to conn if server told to do so
			client.writeMessage(msg)

		case msg := <-client.readChan:
			// Handle messages that were read in readPump
			client.handleReadMessage(msg, roomRT)

		case <-client.flow.Stopped():
			// when server asked to stop running
			return
		}
	}
}

// constantly read messages from connection
func (client *RoomClient) readPump() {
	log.Println("<RoomClient ReadPump>")

	defer func() {
		client.flow.Stop()
		log.Println("<RoomClient ReadPump End>")
	}()

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
			// unmarshaling error
			continue

		default:
			// reading error
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("error reading ws message:", err)
			}
			return
		}

		select {
		case client.readChan <- msg:
			// send message to read chan
		case <-client.flow.Done():
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
		client.stopWithMessage("Unexpected error!")
	}
}

// write message to connection
func (client *RoomClient) writeMessage(msg *message.JSON) {
	client.conn.SetWriteDeadline(time.Now().Add(writeWait)) // if WriteMessage can't send message in writeWait period, client is disconnected

	err := client.conn.WriteJSON(msg)
	// if couldn't write message - disconnect
	if err != nil {
		client.flow.Stop()
	}
}

// process read message
func (client *RoomClient) handleReadMessage(msg *message.JSON, roomRT *RoomRT) {
	// simply tell room about read message
	roomRT.readChan <- RoomReadMessage{
		client:  client,
		message: msg,
	}
}

// send message to client and close after
func (client *RoomClient) stopWithMessage(text string) {
	client.writeChan <- &message.JSON{
		Type: "message",
		Body: text,
	}
	client.flow.Stop()
}
