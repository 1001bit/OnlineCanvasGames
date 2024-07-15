package clicker

import (
	"github.com/1001bit/OnlineCanvasGames/internal/gamelogic"
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
	rtclient "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/client"
)

type ClickerGL struct {
	clicks     uint
	maxPlayers int
}

func NewClickerGL() *ClickerGL {
	return &ClickerGL{
		clicks:     0,
		maxPlayers: 10,
	}
}

func (gl *ClickerGL) Run(doneChan <-chan struct{}, writer gamelogic.RoomWriter) {

}

func (gl *ClickerGL) HandleReadMessage(msg rtclient.MessageWithClient, writer gamelogic.RoomWriter) {
	if msg.Message.Type == "click" {
		gl.clicks += 1
		writer.GlobalWriteMessage(gl.newStateMessage())
	}
}

func (gl *ClickerGL) JoinClient(userID int, writer gamelogic.RoomWriter) {
	writer.WriteMessageTo(gl.newStateMessage(), userID)
}

func (gl *ClickerGL) DeleteClient(userID int, writer gamelogic.RoomWriter) {

}

func (gl *ClickerGL) GetMaxClients() int {
	return gl.maxPlayers
}

func (gl *ClickerGL) newStateMessage() *message.JSON {
	return &message.JSON{
		Type: "clicks",
		Body: gl.clicks,
	}
}
