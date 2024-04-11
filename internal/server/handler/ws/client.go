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
	conn *websocket.Conn

	room *GameRoom

	writeChan chan string
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,

		room: nil,

		writeChan: make(chan string),
	}
}

func (c *Client) closeConn() {
	c.conn.WriteMessage(websocket.CloseMessage, []byte("closed!"))
	c.conn.Close()
}

func (c *Client) readPump() {
	defer func() {
		c.room.disconnectChan <- c
		c.closeConn()
	}()

	// On Pong
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		// Get message from client and send it to hub
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("error reading ws message:", err)
			}
			break
		}

		c.room.readChan <- ClientMessage{
			client: c,
			text:   string(message),
		}
	}
}

func (c *Client) writePump() {
	log.Println("WS Client WritePump")

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.closeConn()
		log.Println("WS Client WritePump End")
	}()

	for {
		select {
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// Ping
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case message, ok := <-c.writeChan:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			// if hub closed c.writeChan
			if !ok {
				return
			}

			c.conn.WriteMessage(websocket.TextMessage, []byte(message))

			log.Println("<Client Write>:", string(message))
		}
	}
}
