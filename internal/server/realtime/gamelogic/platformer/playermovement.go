package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/physics"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/gamelogic"
)

const playerSpeed = 0.01

func (l *Level) ControlPlayerRect(input gamelogic.UserInput) {
	rectID := l.playersRects[input.UserID]
	kRect, ok := l.physEnv.GetKinematicRect(rectID)
	if !ok {
		return
	}

	add := physics.NewVector2(0, 0)
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
