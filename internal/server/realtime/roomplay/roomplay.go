package roomplay

import (
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/rtclient"
)

type RoomWriter interface {
	GlobalWriteMessage(msg *message.JSON)
	WriteMessageTo(msg *message.JSON, id int)
}

type RoomPlay interface {
	Run(doneChan <-chan struct{}, writer RoomWriter)
	HandleReadMessage(msg rtclient.MessageWithClient)
}
