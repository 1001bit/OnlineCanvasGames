package platformer

import (
	"math"

	"github.com/1001bit/OnlineCanvasGames/internal/gamelogic"
)

// InputMap with coeffs instead of ticks
type PlayerInput struct {
	gamelogic.InputMap

	// serverTPS/clientTPS
	serverClientTpsRatio float64
}

func NewPlayerInput(inputMap gamelogic.InputMap, serverTPS, clientTPS float64) *PlayerInput {
	return &PlayerInput{
		InputMap: inputMap,

		serverClientTpsRatio: serverTPS / clientTPS,
	}
}

func (inputMap *PlayerInput) GetHoldCoeff(control string) (float64, bool) {
	ticks, ok := inputMap.GetTicks(control)
	if !ok || ticks == 0 {
		return 0, false
	}

	// max ticks = ceil(clientTPS / serverTPS). Basically, how many times client can tick before server tick
	maxTicks := int(math.Ceil(1 / inputMap.serverClientTpsRatio))
	ticks = min(ticks, maxTicks)

	return float64(ticks) * inputMap.serverClientTpsRatio, true
}

// Add input to map
func (l *Level) HandleInput(userID int, inputMap gamelogic.InputMap) {
	l.userInputMutex.Lock()
	defer l.userInputMutex.Unlock()

	l.userInput[userID] = NewPlayerInput(inputMap, l.tps, l.clientTPS)
}

func (l *Level) ClearUserInput() {
	clear(l.userInput)
}
