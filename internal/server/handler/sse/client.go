package sse

import (
	"fmt"
	"log"
	"net/http"
)

// Client
type Client struct {
	writer http.ResponseWriter

	hub *GameHub

	writeChan chan string
}

func NewClient(writer http.ResponseWriter) *Client {
	return &Client{
		writer: writer,

		hub: nil,

		writeChan: make(chan string),
	}
}

func (c *Client) writePump(done <-chan struct{}) {
	log.Println("<SSE Client WritePump>")

	defer func() {
		c.hub.disconnectClientChan <- c
		log.Println("<SSE Client WritePump End>")
	}()

	for {
		select {
		case message, ok := <-c.writeChan:
			// if hub closed c.write chan
			if !ok {
				return
			}

			c.writeMessage(message)
			log.Println("<Client Write>:", string(message))

		case <-done:
			return
		}
	}
}

func (c *Client) writeMessage(message string) {
	fmt.Fprintf(c.writer, "data: %s\n\n", message)
	c.writer.(http.Flusher).Flush()
}
