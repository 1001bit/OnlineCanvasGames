package platformer

import (
	rtclient "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/client"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/gamelogic"
)

type PlatformerGL struct {
}

func NewPlatformerRP() *PlatformerGL {
	return &PlatformerGL{}
}

func (gl *PlatformerGL) Run(doneChan <-chan struct{}, writer gamelogic.RoomWriter) {

}

func (gl *PlatformerGL) HandleReadMessage(msg rtclient.MessageWithClient, writer gamelogic.RoomWriter) {

}

func (gl *PlatformerGL) JoinClient(userID int, writer gamelogic.RoomWriter) {

}

func (gl *PlatformerGL) GetMaxClients() int {
	return 0
}
