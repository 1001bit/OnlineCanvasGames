package clicker

import (
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
	rtclient "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/client"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/roomplay"
)

type ClickerRP struct {
	clicks uint
}

func NewClickerRP() *ClickerRP {
	return &ClickerRP{
		clicks: 0,
	}
}

func (rp *ClickerRP) Run(doneChan <-chan struct{}, writer roomplay.RoomWriter) {

}

func (rp *ClickerRP) HandleReadMessage(msg rtclient.MessageWithClient, writer roomplay.RoomWriter) {
	if msg.Message.Type == "click" {
		rp.clicks += 1
		writer.GlobalWriteMessage(rp.newStateMessage())
	}
}

func (rp *ClickerRP) JoinClient(userID int, writer roomplay.RoomWriter) {
	writer.WriteMessageTo(rp.newStateMessage(), userID)
}

func (rp *ClickerRP) GetMaxClients() int {
	return 10
}

func (rp *ClickerRP) newStateMessage() *message.JSON {
	return &message.JSON{
		Type: "clicks",
		Body: rp.clicks,
	}
}
