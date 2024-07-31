package clicker

import (
	"github.com/neinBit/ocg-games-service/internal/gamelogic"
	"github.com/neinBit/ocg-games-service/internal/server/message"
	rtclient "github.com/neinBit/ocg-games-service/internal/server/realtime/client"
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

func (gl *ClickerGL) JoinClient(username string, writer gamelogic.RoomWriter) {
	writer.WriteMessageTo(gl.newStateMessage(), username)
}

func (gl *ClickerGL) DeleteClient(username string, writer gamelogic.RoomWriter) {

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
