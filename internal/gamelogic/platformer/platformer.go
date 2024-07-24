package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/gamelogic"
	rtclient "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/client"
	"github.com/1001bit/OnlineCanvasGames/pkg/gameloop"
)

type PlatformerGL struct {
	level *Level

	maxPlayers int
	tps        int
}

func NewPlatformerGL() *PlatformerGL {
	return &PlatformerGL{
		level: NewPlatformerLevel(),

		maxPlayers: 8,
		tps:        25,
	}
}

func (gl *PlatformerGL) Run(doneChan <-chan struct{}, writer gamelogic.RoomWriter) {
	go gameloop.Gameloop(func(dtMs float64) {
		gl.tick(dtMs, writer)
	}, gl.tps, doneChan)
}

func (gl *PlatformerGL) HandleReadMessage(msg rtclient.MessageWithClient, writer gamelogic.RoomWriter) {
	// TODO: Handle Input
}

func (gl *PlatformerGL) JoinClient(userID int, writer gamelogic.RoomWriter) {
	rectID, rect := gl.level.CreatePlayer(userID, gl.maxPlayers)

	writer.WriteMessageTo(NewLevelMessage(gl.level, rectID), userID)

	writer.GlobalWriteMessage(NewConnectMessage(rectID, rect))
}

func (gl *PlatformerGL) DeleteClient(userID int, writer gamelogic.RoomWriter) {
	rectID, err := gl.level.DeletePlayer(userID)
	if err == nil {
		writer.GlobalWriteMessage(NewDisconnectMessage(rectID))
	}
}

func (gl *PlatformerGL) GetMaxClients() int {
	return gl.maxPlayers
}
