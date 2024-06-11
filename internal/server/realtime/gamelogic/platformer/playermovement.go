package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/physics"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/gamelogic"
)

const playerSpeed = 1.0 / 1000

func (l *Level) ControlPlayerRect(input gamelogic.UserInput) {
	rectID := l.playersRects[input.UserID]
	kRect, ok := l.physEnv.GetKinematicRects()[rectID]
	if !ok {
		return
	}

	add := physics.Vector2f{
		X: 0,
		Y: 0,
	}
	if control, ok := input.Controls["left"]; ok {
		add.X -= playerSpeed * control.HoldPeriod
	}
	if control, ok := input.Controls["right"]; ok {
		add.X += playerSpeed * control.HoldPeriod
	}
	if control, ok := input.Controls["up"]; ok {
		add.Y -= playerSpeed * control.HoldPeriod
	}
	if control, ok := input.Controls["down"]; ok {
		add.Y += playerSpeed * control.HoldPeriod
	}

	kRect.AddToAccel(add)
}
