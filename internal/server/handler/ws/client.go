package ws

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

// Client
type Client struct {
	conn      *websocket.Conn
	room      *GameRoom
	userID    int
	writeChan chan string
}

func NewClient(conn *websocket.Conn, userID int) *Client {
	return &Client{
		conn:      conn,
		room:      nil,
		userID:    userID,
		writeChan: make(chan string),
	}
}

// close websocket connection
func (c *Client) closeConn() {
	c.conn.WriteMessage(websocket.CloseMessage, []byte("closed!"))
	c.conn.Close()
}

// constantly read messages from connection
func (c *Client) readPump() {
	log.Println("<WS Client ReadPump>")

	defer func() {
		c.room.disconnectClientChan <- c
		c.closeConn()
		log.Println("<WS Client ReadPump End>")
	}()

	// On Pong
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		// Get message from client
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("error reading ws message:", err)
			}
			break
		}

		// send the message to hub
		c.room.readChan <- ClientMessage{
			client: c,
			text:   string(message),
		}
	}
}

// constantly check messages in writeChan and send them to connection
func (c *Client) writePump() {
	log.Println("<WS Client WritePump>")

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.closeConn()
		log.Println("<WS Client WritePump End>")
	}()

	for {
		select {
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// Ping on ticker clear
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case message, ok := <-c.writeChan:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			// if hub closed c.writeChan
			if !ok {
				return
			}

			// write message from writeChan to connection
			c.conn.WriteMessage(websocket.TextMessage, []byte(message))

			log.Println("<Client Write>:", string(message))
		}
	}
}
