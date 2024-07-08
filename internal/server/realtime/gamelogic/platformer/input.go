package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/physics"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/gamelogic"
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
	playerKRect, ok := l.physEnv.GetKinematicRects()[rectID]
	if !ok {
		return
	}

	add := physics.Vector2f{
		X: 0,
		Y: 0,
	}

	if input.IsControlHeld("left") {
		add.X -= platformerConstants.PlayerSpeed
	}
	if input.IsControlHeld("right") {
		add.X += platformerConstants.PlayerSpeed
	}
	if input.IsControlHeld("jump") && playerKRect.GetCollisionDir().Vertical == physics.Down {
		add.Y -= platformerConstants.PlayerJump
	}

	playerKRect.AddToVel(add)
}
