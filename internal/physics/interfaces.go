package physics

import (
	"github.com/1001bit/OnlineCanvasGames/internal/mathobjects"
)

type Rect interface {
	SetPosition(x, y float64)
	GetPosition() mathobjects.Vector2[float64]
	GetSize() mathobjects.Vector2[float64]
}

type Physical interface {
	Rect
	GetRect() mathobjects.Rect
	CanCollide() bool
}

type Kinematic interface {
	Physical
	AddToVel(x, y float64)
	GetVelocity() mathobjects.Vector2[float64]
	SetVelocity(x, y float64)
	SetCollisionDir(dir mathobjects.Direction)
	DoApplyForce(force ForceType) bool
}
