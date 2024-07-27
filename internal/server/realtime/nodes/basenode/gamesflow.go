package basenode

import (
	"log"

	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/gamenode"
)

func (baseNode *BaseNode) gamesFlow() {
	log.Println("--<BaseNode gamesFlow>")
	defer log.Println("--<BaseNode gamesFlow Done>")

	for {
		select {
		case game := <-baseNode.games.ToConnect():
			// When server asked to connect new game
			baseNode.connectGame(game)
			log.Println("<BaseNode +Game>:", baseNode.games.IDMap.Length())

		case game := <-baseNode.games.ToDisconnect():
			// When server asked to disconnect a game
			baseNode.disconnectGame(game)
			log.Println("<BaseNode -Game>:", baseNode.games.IDMap.Length())

		case <-baseNode.Flow.Done():
			return
		}
	}
}

// connect gameNode to BaseNode
func (baseNode *BaseNode) connectGame(game *gamenode.GameNode) {
	baseNode.games.IDMap.Set(game.GetGame().ID, game)
}

// disconnect gameNode from BaseNode
func (baseNode *BaseNode) disconnectGame(game *gamenode.GameNode) {
	baseNode.games.IDMap.Delete(game.GetGame().ID)
}
