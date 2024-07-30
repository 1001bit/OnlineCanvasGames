package roomnode

import (
	"log"
	"time"

	"github.com/neinBit/ocg-games-service/internal/gamelogic"
	"github.com/neinBit/ocg-games-service/internal/server/realtime/children"
	"github.com/neinBit/ocg-games-service/internal/server/realtime/nodes/roomclient"
	"github.com/neinBit/ocg-games-service/internal/server/realtime/runflow"
)

const roomStopWait = 5 * time.Second

// Layer of RT which is responsible for connections: RoomNode > RoomClients
type RoomNode struct {
	Flow runflow.RunFlow

	Clients children.ChildrenWithID[roomclient.RoomClient]

	connectedToGameChan chan struct{}

	gamelogic gamelogic.GameLogic

	id    int
	owner *roomclient.RoomClient
}

func NewRoomNode(gameID int) *RoomNode {
	return &RoomNode{
		Flow: runflow.MakeRunFlow(),

		Clients: children.MakeChildrenWithID[roomclient.RoomClient](),

		connectedToGameChan: make(chan struct{}),

		gamelogic: NewGameLogicByID(gameID),

		id:    0,
		owner: nil,
	}
}

func (roomNode *RoomNode) Run(requester GameNodeRequester) {
	defer roomNode.Flow.CloseDone()

	log.Println("-<RoomNode Run>")
	defer log.Println("-<RoomNode Run Done>")

	stopTimer := time.NewTimer(roomStopWait)
	defer stopTimer.Stop()

	go roomNode.clientsFlow(requester, stopTimer)
	go roomNode.gamelogic.Run(roomNode.Flow.Done(), roomNode)

	for {
		select {
		case <-stopTimer.C:
			// If timer is over, check for clients
			if roomNode.Clients.IDMap.Length() == 0 {
				return
			}

		case <-roomNode.Flow.Stopped():
			// When server asked to stop running
			return
		}
	}
}
