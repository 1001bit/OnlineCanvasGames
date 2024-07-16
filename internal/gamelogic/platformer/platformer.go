package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/gamelogic"
	rtclient "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/client"
	"github.com/1001bit/OnlineCanvasGames/pkg/gameloop"
)

type PlatformerGL struct {
	level      *Level
	maxPlayers int

	tps int

	inputChan chan gamelogic.UserInput
}

func NewPlatformerGL() *PlatformerGL {
	const (
		maxPlayers = 2
		tps        = 50
	)

	return &PlatformerGL{
		level:      NewPlatformerLevel(),
		maxPlayers: maxPlayers,

		tps: tps,

		inputChan: make(chan gamelogic.UserInput),
	}
}

func (gl *PlatformerGL) Run(doneChan <-chan struct{}, writer gamelogic.RoomWriter) {
	go gameloop.Gameloop(func(dtMs float64) {
		gl.tick(dtMs, writer)
	}, gl.tps, doneChan)

	<-doneChan
}

func (gl *PlatformerGL) HandleReadMessage(msg rtclient.MessageWithClient, writer gamelogic.RoomWriter) {
	switch msg.Message.Type {
	case "input":
		gamelogic.ExtractInputFromMsg(msg.Message.Body, msg.Client.GetUser().ID, gl.inputChan)
	}
}

func (gl *PlatformerGL) JoinClient(userID int, writer gamelogic.RoomWriter) {
	rectID := gl.level.CreatePlayer(userID, gl.maxPlayers)

	writer.WriteMessageTo(gl.NewGameInfoMessage(rectID, gl.tps, platformerConstants), userID)
	writer.WriteMessageTo(gl.NewLevelMessage(), userID)

	rect := gl.level.physEng.GetKinematicRects()[rectID]

	writer.GlobalWriteMessage(gl.NewCreateMessage(rectID, rect))
}

func (gl *PlatformerGL) DeleteClient(userID int, writer gamelogic.RoomWriter) {
	rectID, err := gl.level.DeletePlayer(userID)
	if err == nil {
		writer.GlobalWriteMessage(gl.NewDeleteMessage(rectID))
	}
}

func (gl *PlatformerGL) GetMaxClients() int {
	return gl.maxPlayers
}
