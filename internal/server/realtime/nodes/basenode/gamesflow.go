package basenode

import (
	"log"

	gamenode "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/gamenode"
)

func (baseNode *BaseNode) gamesFlow() {
	log.Println("--<BaseNode gamesFlow>")
	defer log.Println("--<BaseNode gamesFlow Done>")

	for {
		select {
		case game := <-baseNode.games.ToConnect():
			// When server asked to connect new game
			baseNode.connectGame(game)
			log.Println("<BaseNode +Game>:", len(baseNode.games.IDMap))

		case game := <-baseNode.games.ToDisconnect():
			// When server asked to disconnect a game
			baseNode.disconnectGame(game)
			log.Println("<BaseNode -Game>:", len(baseNode.games.IDMap))

		case <-baseNode.Flow.Done():
			return
		}
	}
}

// connect gameNode to BaseNode
func (baseNode *BaseNode) connectGame(game *gamenode.GameNode) {
	baseNode.games.IDMap[game.GetGame().ID] = game
}

// disconnect gameNode from BaseNode
func (baseNode *BaseNode) disconnectGame(game *gamenode.GameNode) {
	delete(baseNode.games.IDMap, game.GetGame().ID)
}
