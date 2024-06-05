package platformer

import (
	rtclient "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/client"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/roomplay"
)

type PlatformerRP struct {
}

func NewPlatformerRP() *PlatformerRP {
	return &PlatformerRP{}
}

func (rp *PlatformerRP) Run(doneChan <-chan struct{}, writer roomplay.RoomWriter) {

}

func (rp *PlatformerRP) HandleReadMessage(msg rtclient.MessageWithClient, writer roomplay.RoomWriter) {

}

func (rp *PlatformerRP) JoinClient(userID int, writer roomplay.RoomWriter) {

}

func (rp *PlatformerRP) GetMaxClients() int {
	return 0
}
