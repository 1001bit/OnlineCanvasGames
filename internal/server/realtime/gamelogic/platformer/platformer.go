package platformer

import (
	rtclient "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/client"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/gamelogic"
)

type PlatformerGL struct {
	level      *Level
	tps        int
	maxPlayers int

	inputChan chan gamelogic.UserInput
}

func NewPlatformerGL() *PlatformerGL {
	return &PlatformerGL{
		level:      NewPlatformerLevel(),
		tps:        20,
		maxPlayers: 2,

		inputChan: make(chan gamelogic.UserInput),
	}
}

func (gl *PlatformerGL) Run(doneChan <-chan struct{}, writer gamelogic.RoomWriter) {
	go gl.gameLoop(doneChan, writer)

	<-doneChan
}

func (gl *PlatformerGL) HandleReadMessage(msg rtclient.MessageWithClient, writer gamelogic.RoomWriter) {
	switch msg.Message.Type {
	case "input":
		gamelogic.ExtractInputFromMsg(msg.Message.Body, msg.Client.GetUser().ID, gl.inputChan)
	}
}

func (gl *PlatformerGL) JoinClient(userID int, writer gamelogic.RoomWriter) {
	rectID := gl.level.CreatePlayer(userID)

	writer.WriteMessageTo(gl.NewGameInfoMessage(rectID), userID)
	writer.WriteMessageTo(gl.NewFullLevelMessage(), userID)
}

func (gl *PlatformerGL) DeleteClient(userID int, writer gamelogic.RoomWriter) {
	gl.level.DeletePlayer(userID)
}

func (gl *PlatformerGL) GetMaxClients() int {
	return gl.maxPlayers
}
