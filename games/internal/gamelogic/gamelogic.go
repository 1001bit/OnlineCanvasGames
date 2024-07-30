package gamelogic

import (
	"github.com/neinBit/ocg-games-service/internal/server/message"
	rtclient "github.com/neinBit/ocg-games-service/internal/server/realtime/client"
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
