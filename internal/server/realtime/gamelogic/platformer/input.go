package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/physics"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/gamelogic"
)

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
	const (
		playerSpeed = 5
		jumpForce   = 3
	)

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
	if input.IsControlHeld("jump") && playerKRect.GetCollisionDir().Vertical == physics.Down {
		add.Y -= jumpForce
	}

	playerKRect.AddToVel(add)
}
