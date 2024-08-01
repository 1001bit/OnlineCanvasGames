package roomclient

import (
	"log"
	"time"

	"github.com/1001bit/ocg-games-service/internal/server/message"
	rtclient "github.com/1001bit/ocg-games-service/internal/server/realtime/client"
	"github.com/1001bit/ocg-games-service/internal/server/realtime/runflow"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
	pongWait  = 60 * time.Second
	// must be shorter that pong wait, because ping should be sent before, than read deadline happens
	pingPeriod = pongWait * 9 / 10

	CloseMsgType = "close"
)

type RoomNodeReader interface {
	ReadMessage(message rtclient.MessageWithClient)
}

// Layer of RT which is responsible for handling connections: RoomClient > User, Room > Roomclient
type RoomClient struct {
	Flow runflow.RunFlow

	conn *websocket.Conn
	user rtclient.User

	writeChan chan *message.JSON
}

func NewRoomClient(conn *websocket.Conn, user rtclient.User) *RoomClient {
	return &RoomClient{
		Flow: runflow.MakeRunFlow(),

		conn: conn,
		user: user,

		writeChan: make(chan *message.JSON),
	}
}

func (client *RoomClient) Run(roomNodeReader RoomNodeReader) {
	defer client.Flow.CloseDone()
	defer client.conn.Close()

	log.Println("-<RoomClient Run>")
	defer log.Println("-<RoomClient Run Done>")

	go client.readFlow(roomNodeReader)
	go client.writeFlow()

	<-client.Flow.Stopped()
}
