package gamelogic

import (
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
	rtclient "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/client"
)

type RoomWriter interface {
	GlobalWriteMessage(msg *message.JSON)
	WriteMessageTo(msg *message.JSON, id int)
}

type GameLogic interface {
	Run(doneChan <-chan struct{}, writer RoomWriter)
	HandleReadMessage(msg rtclient.MessageWithClient, writer RoomWriter)
	JoinClient(userID int, writer RoomWriter)
	DeleteClient(userID int, writer RoomWriter)
	GetMaxClients() int
}