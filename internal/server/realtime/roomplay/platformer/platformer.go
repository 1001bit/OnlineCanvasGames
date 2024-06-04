package platformer

import (
	rtclient "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/client"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/roomplay"
)

type PlatformerRP struct {
	clientChan chan int
}

func NewPlatformerRP() *PlatformerRP {
	return &PlatformerRP{
		clientChan: make(chan int),
	}
}

func (rp *PlatformerRP) Run(doneChan <-chan struct{}, writer roomplay.RoomWriter) {
	for {

	}
}

func (rp *PlatformerRP) HandleReadMessage(msg rtclient.MessageWithClient) {

}

func (rp *PlatformerRP) JoinClient(userID int) {
	rp.clientChan <- userID
}
