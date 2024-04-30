package basenode

import (
	"context"
	"log"

	gamemodel "github.com/1001bit/OnlineCanvasGames/internal/model/game"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/children"
	gamenode "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/gamenode"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/roomclient"
)

// Basic layer of RT which is responsible for handling Games and room-client connections
type BaseNode struct {
	games children.ChildrenWithID[gamenode.GameNode]

	roomsClients children.ChildrenWithID[roomclient.RoomClient]

	gamesJSON []gamemodel.Game
}

func NewBaseNode() *BaseNode {
	return &BaseNode{
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
	log.Println("<BaseNode Run>")
	defer log.Println("<BaseNode Done>")

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

		case client := <-baseNode.roomsClients.ToConnect():
			// When new WS connection needs to be created
			baseNode.protectRoomClient(client)

		case client := <-baseNode.roomsClients.ToDisconnect():
			// When client is done
			baseNode.deleteRoomClient(client)
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

// if there is already client with such ID - stop them. Put a new one
func (baseNode *BaseNode) protectRoomClient(client *roomclient.RoomClient) {
	oldClient, ok := baseNode.roomsClients.IDMap[client.GetID()]
	if ok {
		oldClient.StopWithMessage("This user has just joined another room")
	}
	baseNode.roomsClients.IDMap[client.GetID()] = client
}

// if there is already client with such ID - stop them. Put a new one
func (baseNode *BaseNode) deleteRoomClient(client *roomclient.RoomClient) {
	// can delete only exact client, not just with the same id
	if baseNode.roomsClients.IDMap[client.GetID()] != client {
		return
	}

	delete(baseNode.roomsClients.IDMap, client.GetID())
}
