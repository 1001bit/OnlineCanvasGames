package platformer

import (
	"encoding/json"

	"github.com/neinBit/ocg-games-service/internal/gamelogic"
	rtclient "github.com/neinBit/ocg-games-service/internal/server/realtime/client"
	"github.com/neinBit/ocg-games-service/pkg/concurrent"
	"github.com/neinBit/ocg-games-service/pkg/gameloop"
)

type PlatformerGL struct {
	level *Level

	// set[userID] of already read clients
	handledClients concurrent.ConcurrentSet[int]

	maxPlayers int
}

func NewPlatformerGL() *PlatformerGL {
	return &PlatformerGL{
		level: NewPlatformerLevel(),

		handledClients: concurrent.MakeSet[int](),

		maxPlayers: 8,
	}
}

func (gl *PlatformerGL) Run(doneChan <-chan struct{}, writer gamelogic.RoomWriter) {
	go gameloop.Gameloop(func(dtMs float64) {
		gl.level.Tick(dtMs, writer)

		// clear read clients set (ready to read new messages)
		gl.handledClients.Clear()
	}, int(gl.level.serverTPS), doneChan)
}

func (gl *PlatformerGL) HandleReadMessage(msg rtclient.MessageWithClient, writer gamelogic.RoomWriter) {
	// Protect from reading the same client many times
	if gl.handledClients.Has(msg.Client.GetUser().ID) {
		return
	}
	gl.handledClients.Insert(msg.Client.GetUser().ID)

	switch msg.Message.Type {
	case "input":
		var inputMap gamelogic.InputMap
		err := json.Unmarshal([]byte(msg.Message.Body.(string)), &inputMap)
		if err != nil {
			return
		}

		gl.level.HandleInput(msg.Client.GetUser().ID, inputMap)
	}
}

func (gl *PlatformerGL) JoinClient(userID int, writer gamelogic.RoomWriter) {
	// create player on server
	rectID, rect := gl.level.CreatePlayer(userID, gl.maxPlayers)
	// write level message to new player
	writer.WriteMessageTo(NewLevelMessage(gl.level, rectID), userID)
	// write new player message to everybody
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