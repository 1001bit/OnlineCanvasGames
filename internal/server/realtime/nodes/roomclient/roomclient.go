package roomclient

import (
	"io"
	"log"
	"time"

	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/rtclient"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/runflow"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
	pongWait  = 60 * time.Second
	// must be shorter that pong wait, because ping should be sent before, than read deadline happens
	pingPeriod = pongWait * 9 / 10
)

type RoomNodeReader interface {
	ReadMessage(message rtclient.MessageWithClient)
}

// Layer of RT which is responsible for handling connection WS
type RoomClient struct {
	Flow runflow.RunFlow

	conn *websocket.Conn
	user rtclient.User

	writeChan chan *message.JSON
	readChan  chan *message.JSON
}

func NewRoomClient(conn *websocket.Conn, user rtclient.User) *RoomClient {
	return &RoomClient{
		Flow: runflow.MakeRunFlow(),

		conn: conn,
		user: user,

		writeChan: make(chan *message.JSON),
		readChan:  make(chan *message.JSON),
	}
}

func (client *RoomClient) Run(roomNodeReader RoomNodeReader) {
	log.Println("<RoomClient Run>")

	// ticker that indicates the need to send ping message
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		client.conn.Close()
		client.Flow.CloseDone()

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
			client.writeMessageToConn(msg)

		case msg := <-client.readChan:
			// Handle messages that were read in readPump
			client.handleReadMessage(msg, roomNodeReader)

		case <-client.Flow.Stopped():
			// when server asked to stop running
			return
		}
	}
}

// constantly read messages from connection
func (client *RoomClient) readPump() {
	log.Println("<RoomClient ReadPump>")

	defer func() {
		go client.Flow.Stop()
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
		case <-client.Flow.Done():
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
		client.StopWithMessage("Unexpected error!")
	}
}

// write message to connection
func (client *RoomClient) writeMessageToConn(msg *message.JSON) {
	client.conn.SetWriteDeadline(time.Now().Add(writeWait)) // if WriteMessage can't send message in writeWait period, client is disconnected

	err := client.conn.WriteJSON(msg)
	// if couldn't write message - disconnect
	if err != nil {
		go client.Flow.Stop()
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
