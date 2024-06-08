package platformer

import (
	rtclient "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/client"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/gamelogic"
)

type PlatformerGL struct {
	level          *Level
	ticksPerSecond int
}

type GameInfo struct {
	TPS int `json:"tps"`
}

func NewPlatformerGL() *PlatformerGL {
	return &PlatformerGL{
		level:          NewPlatformerLevel(),
		ticksPerSecond: 20,
	}
}

func (gl *PlatformerGL) Run(doneChan <-chan struct{}, writer gamelogic.RoomWriter) {
	go gl.gameLoop(doneChan)

	<-doneChan
}

func (gl *PlatformerGL) HandleReadMessage(msg rtclient.MessageWithClient, writer gamelogic.RoomWriter) {
	// TODO: Handle Input
}

func (gl *PlatformerGL) JoinClient(userID int, writer gamelogic.RoomWriter) {
	gameinfo := GameInfo{
		TPS: gl.ticksPerSecond,
	}

	writer.WriteMessageTo(gamelogic.NewGameInfoMessage(gameinfo), userID)
	writer.WriteMessageTo(NewLevelMessage(gl.level), userID)
}

func (gl *PlatformerGL) GetMaxClients() int {
	return 2
}
