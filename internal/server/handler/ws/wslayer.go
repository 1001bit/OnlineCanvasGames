package ws

type ClientMessage struct {
	client *Client
	text   string
}

type WSLayer struct {
	clients           map[*Client]bool
	connectChan       chan *Client
	disconnectChan    chan *Client
	clientMessageChan chan ClientMessage
}

func MakeWSLayer() WSLayer {
	return WSLayer{
		clients:           make(map[*Client]bool),
		connectChan:       make(chan *Client),
		disconnectChan:    make(chan *Client),
		clientMessageChan: make(chan ClientMessage),
	}
}
