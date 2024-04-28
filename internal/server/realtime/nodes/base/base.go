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
type BaseRT struct {
	games children.ChildrenWithID[gamenode.GameRT]
}

func NewBaseRT() *BaseRT {
	return &BaseRT{
		games: children.MakeChildrenWithID[gamenode.GameRT](),
	}
}

// get all the games from database and put then into BaseRT
func (baseRT *BaseRT) InitGames() error {
	games, err := gamemodel.GetAll(context.Background())
	if err != nil {
		return err
	}

	for _, game := range games {
		gameRT := gamenode.NewGameRT(game.ID)

		// RUN gameRT
		go func() {
			baseRT.games.ConnectChild(gameRT)
			gameRT.Run()
			baseRT.games.DisconnectChild(gameRT)
		}()
	}

	return nil
}

func (baseRT *BaseRT) Run() {
	log.Println("<BaseRT Run>")
	defer log.Println("<BaseRT Done>")

	for {
		select {
		// Games
		case game := <-baseRT.games.ToConnect():
			baseRT.connectGame(game)
			log.Println("<BaseRT +Game>:", len(baseRT.games.IDMap))

		case game := <-baseRT.games.ToDisconnect():
			baseRT.disconnectGame(game)
			log.Println("<BaseRT -Game>:", len(baseRT.games.IDMap))
		}
	}
}

// create new room and connect it to BaseRT
func (baseRT *BaseRT) ConnectNewRoom(ctx context.Context, gameID int) (*roomnode.RoomRT, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	gameRT, ok := baseRT.games.IDMap[gameID]
	if !ok {
		return nil, ErrNoGame
	}

	room := roomnode.NewRoomRT()

	// RUN roomRT
	go func() {
		gameRT.Rooms.ConnectChild(room)
		room.Run(gameRT.RequestUpdatingRoomsJSON)
		gameRT.Rooms.DisconnectChild(room)
	}()

	// wait until room connected to BaseRT
	select {
	case <-room.ConnectedToGame():
		return room, nil
	case <-ctx.Done():
		room.Flow.Stop()
		return nil, ErrCreateRoom
	}
}

func (baseRT *BaseRT) GetGameByID(id int) (*gamenode.GameRT, bool) {
	game, ok := baseRT.games.IDMap[id]
	return game, ok
}

// connect gameRT to BaseRT
func (baseRT *BaseRT) connectGame(game *gamenode.GameRT) {
	baseRT.games.IDMap[game.GetID()] = game
}

// disconnect gameRT from BaseRT
func (baseRT *BaseRT) disconnectGame(game *gamenode.GameRT) {
	delete(baseRT.games.IDMap, game.GetID())
}
