package basenode

import (
	"context"
	"log"

	"github.com/neinBit/ocg-games-service/internal/gamemodel"
	"github.com/neinBit/ocg-games-service/internal/server/realtime/children"
	"github.com/neinBit/ocg-games-service/internal/server/realtime/nodes/gamenode"
	"github.com/neinBit/ocg-games-service/internal/server/realtime/nodes/roomclient"
	"github.com/neinBit/ocg-games-service/internal/server/realtime/runflow"
)

// Basic layer of RT which is responsible for handling connections: BaseNode > GameNodes, BaseNode (>) RoomClients
type BaseNode struct {
	Flow runflow.RunFlow

	games        children.MapChildren[string, *gamenode.GameNode]
	roomsClients children.MapChildren[string, *roomclient.RoomClient]

	gamesJSON []gamemodel.Game
}

func NewBaseNode() *BaseNode {
	return &BaseNode{
		Flow: runflow.MakeRunFlow(),

		games:        children.MakeMapChildren[string, *gamenode.GameNode](),
		roomsClients: children.MakeMapChildren[string, *roomclient.RoomClient](),

		gamesJSON: make([]gamemodel.Game, 0),
	}
}

// get all the games from database and put then into BaseNode
func (baseNode *BaseNode) InitGames() error {
	var err error

	baseNode.gamesJSON, err = gamemodel.GetAll(context.Background())
	if err != nil {
		return err
	}

	for i := range baseNode.gamesJSON {
		gameNode := gamenode.NewGameNode(baseNode.gamesJSON[i])

		// RUN gameNode
		go func() {
			baseNode.games.ConnectChild(gameNode, baseNode.Flow.Done())
			gameNode.Run()
			baseNode.games.DisconnectChild(gameNode, baseNode.Flow.Done())
		}()
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
