package ws

type WSLayer struct {
	clients        map[*Client]bool
	connectChan    chan *Client
	disconnectChan chan *Client
	messageChan    chan string
}

func MakeWSLayer() WSLayer {
	return WSLayer{
		clients:        make(map[*Client]bool),
		connectChan:    make(chan *Client),
		disconnectChan: make(chan *Client),
		messageChan:    make(chan string),
	}
}
