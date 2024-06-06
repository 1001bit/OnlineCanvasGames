package platformer

import (
	rtclient "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/client"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/gamelogic"
)

type PlatformerGL struct {
	level          *gamelogic.Level
	ticksPerSecond int
}

func NewPlatformerGL() *PlatformerGL {
	return &PlatformerGL{
		level:          NewPlatformerLevel(),
		ticksPerSecond: 20,
	}
}

func (gl *PlatformerGL) Run(doneChan <-chan struct{}, writer gamelogic.RoomWriter) {

}

func (gl *PlatformerGL) HandleReadMessage(msg rtclient.MessageWithClient, writer gamelogic.RoomWriter) {

}

func (gl *PlatformerGL) JoinClient(userID int, writer gamelogic.RoomWriter) {
	writer.WriteMessageTo(gamelogic.NewGameDataMessage(gl.ticksPerSecond, gl.level), userID)
}

func (gl *PlatformerGL) GetMaxClients() int {
	return 1
}
