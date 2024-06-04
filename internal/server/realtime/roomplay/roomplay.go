package roomplay

import (
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
	rtclient "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/client"
)

type RoomWriter interface {
	GlobalWriteMessage(msg *message.JSON)
	WriteMessageTo(msg *message.JSON, id int)
}

type RoomPlay interface {
	Run(doneChan <-chan struct{}, writer RoomWriter)
	HandleReadMessage(msg rtclient.MessageWithClient)
	JoinClient(userID int)
	GetMaxClients() int
}
