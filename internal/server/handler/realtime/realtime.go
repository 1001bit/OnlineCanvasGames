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
	games              map[int]*GameRT
	connectGameChan    chan *GameRT
	disconnectGameChan chan *GameRT

	// Not allowing the same user join two rooms at once
	roomsClients             map[int]*RoomClient
	registerRoomClientChan   chan *RoomClient
	unregisterRoomClientChan chan *RoomClient
}

func NewRealtime() *Realtime {
	return &Realtime{
		games:              make(map[int]*GameRT),
		connectGameChan:    make(chan *GameRT),
		disconnectGameChan: make(chan *GameRT),

		roomsClients:             make(map[int]*RoomClient),
		registerRoomClientChan:   make(chan *RoomClient),
		unregisterRoomClientChan: make(chan *RoomClient),
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
		go gameRT.Run()

		go func() {
			rt.connectGameChan <- gameRT
		}()
	}

	return nil
}

func (rt *Realtime) Run() {
	log.Println("<Realtime Run>")
	defer log.Println("<Realtime Done>")

	for {
		select {
		// Games
		case game := <-rt.connectGameChan:
			rt.connectGame(game)
			log.Println("<Realtime +Game>:", len(rt.games))

		case game := <-rt.disconnectGameChan:
			rt.disconnectGame(game)
			log.Println("<Realtime -Game>:", len(rt.games))

		// Rooms clients
		case client := <-rt.registerRoomClientChan:
			rt.registerRoomClient(client)

		case client := <-rt.unregisterRoomClientChan:
			rt.unregisterRoomClient(client)
		}
	}
}

// create new room and connect it to RT
func (rt *Realtime) ConnectNewRoom(ctx context.Context, gameID int) (*RoomRT, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	gameRT, ok := rt.games[gameID]
	if !ok {
		return nil, ErrNoGame
	}

	room := NewRoomRT()
	go room.Run()

	// request connecting room to RT
	select {
	case gameRT.connectRoomChan <- room:
	case <-ctx.Done():
		return nil, ErrCreateRoom
	}

	// wait until room connected to RT
	select {
	case <-room.connectedToGame:
		return room, nil
	case <-ctx.Done():
		return nil, ErrCreateRoom
	}
}

func (rt *Realtime) GetGameByID(id int) (*GameRT, bool) {
	game, ok := rt.games[id]
	return game, ok
}

// connect gameRT to RT
func (rt *Realtime) connectGame(game *GameRT) {
	rt.games[game.gameID] = game
	game.rt = rt
}

// disconnect gameRT from RT
func (rt *Realtime) disconnectGame(game *GameRT) {
	if _, ok := rt.games[game.gameID]; !ok {
		return
	}

	delete(rt.games, game.gameID)
}

// Called by room when a client is connected. Disconnects client with the same id from previous room and puts new into list
func (rt *Realtime) registerRoomClient(client *RoomClient) {
	if oldClient, ok := rt.roomsClients[client.user.id]; ok {
		oldClient.stopWithMessage("This user has just joined another room")
	}

	rt.roomsClients[client.user.id] = client
}

// Called by room when a client is disconnected. Removes client from list if requested client IS the client in the list
func (rt *Realtime) unregisterRoomClient(client *RoomClient) {
	if oldClient, ok := rt.roomsClients[client.user.id]; ok && oldClient == client {
		delete(rt.roomsClients, client.user.id)
	}
}
