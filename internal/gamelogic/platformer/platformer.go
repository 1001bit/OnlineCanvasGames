package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/gamelogic"
	rtclient "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/client"
	"github.com/1001bit/OnlineCanvasGames/pkg/gameloop"
)

type UserInput struct {
	InputMap gamelogic.InputMap
	UserID   int
}

type PlatformerGL struct {
	level *Level

	maxPlayers int
	tps        int

	inputChan chan UserInput
}

func NewPlatformerGL() *PlatformerGL {
	return &PlatformerGL{
		level: NewPlatformerLevel(),

		maxPlayers: 4,
		tps:        25,

		inputChan: make(chan UserInput),
	}
}

func (gl *PlatformerGL) Run(doneChan <-chan struct{}, writer gamelogic.RoomWriter) {
	go gameloop.Gameloop(func(dtMs float64) {
		gl.tick(dtMs, writer)
	}, gl.tps, doneChan)
}

func (gl *PlatformerGL) HandleReadMessage(msg rtclient.MessageWithClient, writer gamelogic.RoomWriter) {
	switch msg.Message.Type {
	case "input":
		inputMap, err := gamelogic.GetInputMapFromMsg(msg.Message.Body, msg.Client.GetUser().ID)
		if err != nil {
			return
		}

		gl.inputChan <- UserInput{
			UserID:   msg.Client.GetUser().ID,
			InputMap: inputMap,
		}
	}
}

func (gl *PlatformerGL) JoinClient(userID int, writer gamelogic.RoomWriter) {
	rectID, rect := gl.level.CreatePlayer(userID, gl.maxPlayers)

	writer.WriteMessageTo(NewGameInfoMessage(gl), userID)
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
