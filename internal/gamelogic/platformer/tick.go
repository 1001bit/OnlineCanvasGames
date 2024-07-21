package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/gamelogic"
)

func (gl *PlatformerGL) tick(dtMs float64, writer gamelogic.RoomWriter) {
	fullInputMap := gl.getFullInputMap()

	moved := gl.level.Tick(dtMs, fullInputMap)

	writer.GlobalWriteMessage(gl.NewUpdateMessage(moved))
}

func (gl *PlatformerGL) getFullInputMap() map[int]gamelogic.InputMap {
	fullInputMap := make(map[int]gamelogic.InputMap)

	select {
	case input := <-gl.inputChan:
		userID := input.UserID
		inputMap := input.InputMap

		// protect from handling input from same user more than once
		if _, ok := fullInputMap[userID]; ok {
			break
		}
		fullInputMap[userID] = inputMap

	default:
	}

	return fullInputMap
}
