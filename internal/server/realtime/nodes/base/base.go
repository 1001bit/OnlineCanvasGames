package basenode

import (
	"context"
	"errors"
	"log"
	"time"

	gamemodel "github.com/1001bit/OnlineCanvasGames/internal/model/game"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/children"
	gamenode "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/game"
	roomnode "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/room"
)

var (
	ErrNoGame     = errors.New("game does not exist")
	ErrNoRoom     = errors.New("room does not exist")
	ErrCreateRoom = errors.New("could not create a room")
)

// Basic layer of RT which is responsible for handling Games and room-client connections
type BaseNode struct {
	games children.ChildrenWithID[gamenode.GameNode]

	roomsClients children.ChildrenWithID[roomnode.RoomClient]
}

func NewBaseNode() *BaseNode {
	return &BaseNode{
		games: children.MakeChildrenWithID[gamenode.GameNode](),

		roomsClients: children.MakeChildrenWithID[roomnode.RoomClient](),
	}
}

// get all the games from database and put then into BaseNode
func (baseNode *BaseNode) InitGames() error {
	games, err := gamemodel.GetAll(context.Background())
	if err != nil {
		return err
	}

	for _, game := range games {
		gameNode := gamenode.NewGameNode(game.ID)

		// RUN gameNode
		go func() {
			baseNode.games.ConnectChild(gameNode)
			gameNode.Run()
			baseNode.games.DisconnectChild(gameNode)
		}()
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

// create new room and connect it to BaseNode
func (baseNode *BaseNode) ConnectNewRoom(ctx context.Context, gameID int) (*roomnode.RoomNode, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	gameNode, ok := baseNode.games.IDMap[gameID]
	if !ok {
		return nil, ErrNoGame
	}

	room := roomnode.NewRoomNode()

	// RUN roomNode
	go func() {
		gameNode.Rooms.ConnectChild(room)
		room.Run(gameNode.RequestUpdatingRoomsJSON)
		gameNode.Rooms.DisconnectChild(room)
	}()

	// wait until room connected to BaseNode
	select {
	case <-room.ConnectedToGame():
		return room, nil
	case <-ctx.Done():
		go room.Flow.Stop()
		return nil, ErrCreateRoom
	}
}

func (baseNode *BaseNode) GetGameByID(id int) (*gamenode.GameNode, bool) {
	game, ok := baseNode.games.IDMap[id]
	return game, ok
}

// connect gameNode to BaseNode
func (baseNode *BaseNode) connectGame(game *gamenode.GameNode) {
	baseNode.games.IDMap[game.GetID()] = game
}

// disconnect gameNode from BaseNode
func (baseNode *BaseNode) disconnectGame(game *gamenode.GameNode) {
	delete(baseNode.games.IDMap, game.GetID())
}

// if there is already client with such ID - stop them. Put a new one
func (baseNode *BaseNode) protectRoomClient(client *roomnode.RoomClient) {
	oldClient, ok := baseNode.roomsClients.IDMap[client.GetID()]
	if ok {
		oldClient.StopWithMessage("This user has just joined another room")
	}
	baseNode.roomsClients.IDMap[client.GetID()] = client
}

// if there is already client with such ID - stop them. Put a new one
func (baseNode *BaseNode) deleteRoomClient(client *roomnode.RoomClient) {
	// can delete only exact client, not just with the same id
	if baseNode.roomsClients.IDMap[client.GetID()] != client {
		return
	}

	delete(baseNode.roomsClients.IDMap, client.GetID())
}
