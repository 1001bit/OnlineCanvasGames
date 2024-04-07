package sse

type SSELayer struct {
	clients     map[*Client]bool
	connect     chan *Client
	disconnect  chan *Client
	messageChan chan string
}

func MakeSSELayer() SSELayer {
	return SSELayer{
		clients:     make(map[*Client]bool),
		connect:     make(chan *Client),
		disconnect:  make(chan *Client),
		messageChan: make(chan string),
	}
}