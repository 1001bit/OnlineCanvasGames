package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/gamelogic"
)

// InputMap with coeffs instead of ticks
type CoeffInputMap struct {
	gamelogic.InputMap

	// serverTPS/clientTPS
	serverClientTpsRatio float64
}

func NewCoeffInputMap(inputMap gamelogic.InputMap, serverTPS, clientTPS float64) *CoeffInputMap {
	return &CoeffInputMap{
		InputMap: inputMap,

		serverClientTpsRatio: serverTPS / clientTPS,
	}
}

func (inputMap *CoeffInputMap) GetHoldCoeff(control string) (float64, bool) {
	ticks, ok := inputMap.GetTicks(control)
	if !ok || ticks == 0 {
		return 0, false
	}

	return float64(ticks) * inputMap.serverClientTpsRatio, true
}

// Add input to map
func (l *Level) HandleInput(userID int, inputMap gamelogic.InputMap) {
	l.userInputMutex.Lock()
	defer l.userInputMutex.Unlock()

	l.userInput[userID] = NewCoeffInputMap(inputMap, l.tps, l.clientTPS)
}

func (l *Level) ClearUserInput() {
	clear(l.userInput)
}
