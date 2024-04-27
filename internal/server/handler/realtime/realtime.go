package realtime

import (
	"context"
	"errors"
	"log"
	"time"

	gamemodel "github.com/1001bit/OnlineCanvasGames/internal/model/game"
)

var (
	ErrNoGame     = errors.New("game does not exist")
	ErrCreateRoom = errors.New("could not create a room")
)

// Basic layer of RT which is responsible for handling Games and room-client connections
type Realtime struct {
	games ChildrenWithID[GameRT]
}

func NewRealtime() *Realtime {
	return &Realtime{
		games: MakeChildrenWithID[GameRT](),
	}
}

// get all the games from database and put then into RT
func (rt *Realtime) InitGames() error {
	games, err := gamemodel.GetAll(context.Background())
	if err != nil {
		return err
	}

	for _, game := range games {
		gameRT := NewGameRT(game.ID)
		go gameRT.Run(rt)
	}

	return nil
}

func (rt *Realtime) Run() {
	log.Println("<Realtime Run>")
	defer log.Println("<Realtime Done>")

	for {
		select {
		// Games
		case game := <-rt.games.connectChan:
			rt.connectGame(game)
			log.Println("<Realtime +Game>:", len(rt.games.idMap))

		case game := <-rt.games.disconnectChan:
			rt.disconnectGame(game)
			log.Println("<Realtime -Game>:", len(rt.games.idMap))
		}
	}
}

// create new room and connect it to RT
func (rt *Realtime) ConnectNewRoom(ctx context.Context, gameID int) (*RoomRT, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	gameRT, ok := rt.games.idMap[gameID]
	if !ok {
		return nil, ErrNoGame
	}

	room := NewRoomRT()
	go room.Run(gameRT)

	// wait until room connected to RT
	select {
	case <-room.connectedToGame:
		return room, nil
	case <-ctx.Done():
		room.flow.Stop()
		return nil, ErrCreateRoom
	}
}

func (rt *Realtime) GetGameByID(id int) (*GameRT, bool) {
	game, ok := rt.games.idMap[id]
	return game, ok
}

// connect gameRT to RT
func (rt *Realtime) connectGame(game *GameRT) {
	rt.games.idMap[game.gameID] = game
}

// disconnect gameRT from RT
func (rt *Realtime) disconnectGame(game *GameRT) {
	delete(rt.games.idMap, game.gameID)
}
