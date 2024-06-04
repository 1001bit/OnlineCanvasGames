package gamenode

import (
	"log"

	gamemodel "github.com/1001bit/OnlineCanvasGames/internal/model/game"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/children"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/gameclient"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/roomnode"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/runflow"
)

// Room that will be sent to client
type RoomJSON struct {
	Owner   string `json:"owner"`
	Clients int    `json:"clients"`
	Limit   int    `json:"limit"`
	ID      int    `json:"id"`
}

// Layer of RT which is responsible for game hub and containing rooms
type GameNode struct {
	Flow runflow.RunFlow

	Rooms   children.ChildrenWithID[roomnode.RoomNode]
	Clients children.Children[gameclient.GameClient]

	roomsJSON           []RoomJSON
	roomsJSONUpdateChan chan struct{}

	game gamemodel.Game
}

func NewGameNode(game gamemodel.Game) *GameNode {
	return &GameNode{
		Flow: runflow.MakeRunFlow(),

		Rooms:   children.MakeChildrenWithID[roomnode.RoomNode](),
		Clients: children.MakeChildren[gameclient.GameClient](),

		roomsJSON:           make([]RoomJSON, 0),
		roomsJSONUpdateChan: make(chan struct{}),

		game: game,
	}
}

func (gameNode *GameNode) Run() {
	defer gameNode.Flow.CloseDone()

	log.Println("-<GameNode Run>")
	defer log.Println("-<GameNode Run Done>")

	go gameNode.clientsFlow()
	go gameNode.roomsFlow()

	<-gameNode.Flow.Stopped()
}
