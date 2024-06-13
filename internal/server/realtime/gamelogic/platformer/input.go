package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/physics"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/gamelogic"
)

const playerSpeed = 5

func (gl *PlatformerGL) handleInput() {
	for {
		select {
		case input := <-gl.inputChan:
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
		add.X -= playerSpeed
	}
	if input.IsControlHeld("right") {
		add.X += playerSpeed
	}
	if input.IsControlHeld("up") {
		add.Y -= playerSpeed
	}
	if input.IsControlHeld("down") {
		add.Y += playerSpeed
	}

	playerKRect.AddToVel(add)
}
