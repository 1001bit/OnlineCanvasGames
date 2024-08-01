package platformer

import (
	"math"

	"github.com/neinBit/ocg-games-service/internal/gamelogic"
)

// InputMap with coeffs instead of ticks
type PlayerInput struct {
	inputMap gamelogic.InputMap

	// serverTPS/clientTPS
	serverClientTpsRatio float64
}

func CreatePlayerInput(serverTPS, clientTPS float64) PlayerInput {
	return PlayerInput{
		inputMap: make(gamelogic.InputMap),

		serverClientTpsRatio: serverTPS / clientTPS,
	}
}

func (input *PlayerInput) GetHoldCoeff(control string) (float64, bool) {
	ticks, ok := input.inputMap.GetTicks(control)
	if !ok || ticks == 0 {
		return 0, false
	}

	// max ticks = ceil(clientTPS / serverTPS). Basically, how many times client can tick before server tick
	maxTicks := int(math.Ceil(1 / input.serverClientTpsRatio))
	ticks = min(ticks, maxTicks)

	return float64(ticks) * input.serverClientTpsRatio, true
}

func (input *PlayerInput) IsHeld(control string) bool {
	return input.inputMap.IsHeld(control)
}

func (input *PlayerInput) SetInputMap(inputMap gamelogic.InputMap) {
	input.inputMap = inputMap
}

func (input *PlayerInput) ClearInputMap() {
	clear(input.inputMap)
}

// Add input to map
func (l *Level) HandleInput(username string, inputMap gamelogic.InputMap) {
	playerData, ok := l.playersData.Get(username)
	if !ok {
		return
	}

	playerData.HandleInput(inputMap)
}
