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
			if _, ok := handledClients[input.UserID]; ok {
				continue
			}
			handledClients[input.UserID] = true

			gl.level.controlPlayerRect(input)
		default:
			return
		}
	}
}

func (l *Level) controlPlayerRect(input gamelogic.UserInput) {
	rectID := l.playersRects[input.UserID]
	playerKRect, ok := l.physEng.GetKinematicRects()[rectID]
	if !ok {
		return
	}

	add := mathobjects.Vector2[float64]{
		X: 0,
		Y: 0,
	}

	if coeff, ok := input.GetControlCoeff("left"); ok {
		add.X -= platformerConstants.PlayerSpeed * coeff
	}

	if coeff, ok := input.GetControlCoeff("right"); ok {
		add.X += platformerConstants.PlayerSpeed * coeff
	}

	if input.IsHeld("jump") && playerKRect.GetCollisionVertical() == mathobjects.Down {
		add.Y -= platformerConstants.PlayerJump
	}

	playerKRect.AddToVel(add)
}
