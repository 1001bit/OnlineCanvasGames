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

	roomsClients             map[int]*RoomRTClient
	connectRoomClientChan    chan *RoomRTClient
	disconnectRoomClientChan chan *RoomRTClient
}

func NewRealtime() *Realtime {
	return &Realtime{
		games:              make(map[int]*GameRT),
		connectGameChan:    make(chan *GameRT),
		disconnectGameChan: make(chan *GameRT),

		roomsClients:             make(map[int]*RoomRTClient),
		connectRoomClientChan:    make(chan *RoomRTClient),
		disconnectRoomClientChan: make(chan *RoomRTClient),
	}
}

// get all the games from database and put then into RT
func (rt *Realtime) InitGames() error {
	games, err := gamemodel.GetAll(context.Background())
	if err != nil {
		return err
	}

	for _, game := range games {
		gameRT := NewGameRT()
		gameRT.gameID = game.ID
		go func() {
			rt.connectGameChan <- gameRT
		}()
	}

	return nil
}

func (rt *Realtime) Run() {
	log.Println("<Realtime Run>")
	defer log.Println("<Realtime Run End>")

	for {
		select {
		case game := <-rt.connectGameChan:
			rt.connectGame(game)
			log.Println("<Realtime +Game>:", len(rt.games))

		case game := <-rt.disconnectGameChan:
			rt.disconnectGame(game)
			log.Println("<Realtime -Game>:", len(rt.games))

		case client := <-rt.connectRoomClientChan:
			rt.connectRoomClient(client)

		case client := <-rt.disconnectRoomClientChan:
			rt.disconnectRoomClient(client)
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

	// request connecting room tw RT
	select {
	case gameRT.connectRoomChan <- room:
	case <-ctx.Done():
		return nil, ErrCreateRoom
	}

	// wait until room connected to RT
	select {
	case <-room.connectedToRT:
		return room, nil
	case <-ctx.Done():
		return nil, ErrCreateRoom
	}

}

// connect gameRT to RT
func (rt *Realtime) connectGame(game *GameRT) {
	rt.games[game.gameID] = game
	game.rt = rt
	go game.Run()
}

// disconnect gameRT from RT
func (rt *Realtime) disconnectGame(game *GameRT) {
	delete(rt.games, game.gameID)
}

// Called by room when a client is connected. Disconnects client with the same id from previous room and puts new into list
func (rt *Realtime) connectRoomClient(client *RoomRTClient) {
	if oldClient, ok := rt.roomsClients[client.userID]; ok {
		oldClient.roomRT.disconnectClientChan <- oldClient
	}

	rt.roomsClients[client.userID] = client
}

// Called by room when a client is disconnected. Removes client from list if requested client IS the client in the list
func (rt *Realtime) disconnectRoomClient(client *RoomRTClient) {
	if oldClient, ok := rt.roomsClients[client.userID]; ok && oldClient == client {
		delete(rt.roomsClients, client.userID)
	}
}
