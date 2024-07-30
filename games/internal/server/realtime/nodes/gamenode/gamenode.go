package gamenode

import (
	"log"

	"github.com/neinBit/ocg-games-service/internal/gamemodel"
	"github.com/neinBit/ocg-games-service/internal/server/realtime/children"
	"github.com/neinBit/ocg-games-service/internal/server/realtime/nodes/gameclient"
	"github.com/neinBit/ocg-games-service/internal/server/realtime/nodes/roomnode"
	"github.com/neinBit/ocg-games-service/internal/server/realtime/runflow"
)

// Room that will be sent to client
type RoomJSON struct {
	Owner   string `json:"owner"`
	Clients int    `json:"clients"`
	Limit   int    `json:"limit"`
	ID      int    `json:"id"`
}

// Layer of RT which is responsible for connections: GameNode > RoomsJson, GameNode > GameClients
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
