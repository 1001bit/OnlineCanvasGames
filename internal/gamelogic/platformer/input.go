package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/gamelogic"
	"github.com/1001bit/OnlineCanvasGames/internal/mathobjects"
)

func (gl *PlatformerGL) handleInput() {
	handledClients := make(map[int]bool)

	for {
		select {
		case input := <-gl.inputChan:
			// protect from handling input from same user more than once
			if _, ok := handledClients[input.GetUserID()]; ok {
				continue
			}
			handledClients[input.GetUserID()] = true

			gl.level.controlPlayerRect(input)

		default:
			return
		}
	}
}

func (l *Level) controlPlayerRect(input *gamelogic.UserInput) {
	rectID := l.playersRects[input.GetUserID()]
	playerKRect, ok := l.physEng.GetKinematicRects()[rectID]
	if !ok {
		return
	}

	addX := 0.0
	addY := 0.0

	if coeff, ok := input.GetControlCoeff("left"); ok {
		addX -= platformerConstants.PlayerSpeed * coeff
	}

	if coeff, ok := input.GetControlCoeff("right"); ok {
		addX += platformerConstants.PlayerSpeed * coeff
	}

	if input.IsHeld("jump") && playerKRect.GetCollisionVertical() == mathobjects.Down {
		addY -= platformerConstants.PlayerJump
	}

	playerKRect.AddToVel(addX, addY)
}
