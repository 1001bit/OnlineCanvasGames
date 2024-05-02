package basenode

import (
	"context"
	"log"

	gamemodel "github.com/1001bit/OnlineCanvasGames/internal/model/game"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/children"
	gamenode "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/gamenode"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/roomclient"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/runflow"
)

// Basic layer of RT which is responsible for handling Games and room-client connections
type BaseNode struct {
	Flow runflow.RunFlow

	games children.ChildrenWithID[gamenode.GameNode]

	roomsClients children.ChildrenWithID[roomclient.RoomClient]

	gamesJSON []gamemodel.Game
}

func NewBaseNode() *BaseNode {
	return &BaseNode{
		Flow: runflow.MakeRunFlow(),

		games: children.MakeChildrenWithID[gamenode.GameNode](),

		roomsClients: children.MakeChildrenWithID[roomclient.RoomClient](),

		gamesJSON: make([]gamemodel.Game, 0),
	}
}

// get all the games from database and put then into BaseNode
func (baseNode *BaseNode) InitGames() error {
	games, err := gamemodel.GetAll(context.Background())
	if err != nil {
		return err
	}

	baseNode.gamesJSON = make([]gamemodel.Game, 0)
	for _, game := range games {
		gameNode := gamenode.NewGameNode(game)

		// RUN gameNode
		go func() {
			baseNode.games.ConnectChild(gameNode)
			gameNode.Run()
			baseNode.games.DisconnectChild(gameNode)
		}()

		// add game to gamesJson
		baseNode.gamesJSON = append(baseNode.gamesJSON, game)
	}

	return nil
}

func (baseNode *BaseNode) Run() {
	defer baseNode.Flow.CloseDone()

	log.Println("-<BaseNode Run>")
	defer log.Println("-<BaseNode Run Done>")

	go baseNode.gamesFlow()
	go baseNode.clientsFlow()

	<-baseNode.Flow.Stopped()
}
