package gameroom

import (
	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/realtime"
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
)

type GameRoom interface {
	// process message that is read from a room's client
	HandleMessage(realtime.RoomReadMessage)
	// Run game room
	Run(doneChan <-chan struct{}, globalWriteChan chan<- *message.JSON)
}
