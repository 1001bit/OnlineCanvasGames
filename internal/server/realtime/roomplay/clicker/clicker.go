package clicker

import (
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
	rtclient "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/client"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/roomplay"
)

type ClickerRP struct {
	clientChan chan int

	clickChan chan struct{}
	clicks    uint
}

func NewClickerRP() *ClickerRP {
	return &ClickerRP{
		clientChan: make(chan int),

		clickChan: make(chan struct{}),
		clicks:    0,
	}
}

func (rp *ClickerRP) Run(doneChan <-chan struct{}, writer roomplay.RoomWriter) {
	for {
		select {
		case <-rp.clickChan:
			writer.GlobalWriteMessage(rp.newStateMessage())
		case userID := <-rp.clientChan:
			writer.WriteMessageTo(rp.newStateMessage(), userID)
		case <-doneChan:
			return
		}
	}
}

func (rp *ClickerRP) HandleReadMessage(msg rtclient.MessageWithClient) {
	if msg.Message.Type == "click" {
		rp.click()
	}
}

func (rp *ClickerRP) JoinClient(userID int) {
	rp.clientChan <- userID
}

func (rp *ClickerRP) newStateMessage() *message.JSON {
	return &message.JSON{
		Type: "clicks",
		Body: rp.clicks,
	}
}

func (rp *ClickerRP) click() {
	rp.clicks += 1
	rp.clickChan <- struct{}{}
}
