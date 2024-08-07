package gamelogic

import (
	rtclient "github.com/1001bit/onlinecanvasgames/services/games/internal/server/realtime/client"
	"github.com/1001bit/onlinecanvasgames/services/games/pkg/message"
)

type RoomWriter interface {
	GlobalWriteMessage(msg *message.JSON)
	WriteMessageTo(msg *message.JSON, name string)
}

type GameLogic interface {
	Run(doneChan <-chan struct{}, writer RoomWriter)
	HandleReadMessage(msg rtclient.MessageWithClient, writer RoomWriter)
	JoinClient(username string, writer RoomWriter)
	DeleteClient(username string, writer RoomWriter)
	GetMaxClients() int
}
