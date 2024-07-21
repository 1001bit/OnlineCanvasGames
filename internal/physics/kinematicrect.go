package physics

import (
	"github.com/1001bit/OnlineCanvasGames/internal/mathobjects"
	"github.com/1001bit/OnlineCanvasGames/pkg/set"
)

// Kinematic Rect
type KinematicRect struct {
	PhysicalRect

	Velocity mathobjects.Vector2[float64] `json:"velocity"`

	ForcesToApply set.Set[ForceType] `json:"forcesToApply"`

	collisionVertical   mathobjects.Direction
	collisionHorizontal mathobjects.Direction
}

func NewKinematicRect(rect PhysicalRect, forces ...ForceType) *KinematicRect {
	return &KinematicRect{
		PhysicalRect: rect,

		Velocity: mathobjects.Vector2[float64]{
			X: 0,
			Y: 0,
		},

		ForcesToApply: set.MakeSet[ForceType](forces),

		collisionVertical:   mathobjects.None,
		collisionHorizontal: mathobjects.None,
	}
}

func (kr *KinematicRect) SetVelocity(x, y float64) {
	kr.Velocity.X = x
	kr.Velocity.Y = y
}

func (kr *KinematicRect) AddToVel(x, y float64) {
	kr.Velocity.X += x
	kr.Velocity.Y += y
}

func (kr *KinematicRect) AddVectorToVel(vector mathobjects.Vector2[float64]) {
	kr.Velocity.Add(vector)
}

func (kr *KinematicRect) SetCollisionDir(dir mathobjects.Direction) {
	if dir == mathobjects.Down || dir == mathobjects.Up {
		kr.collisionVertical = dir
	} else if dir == mathobjects.Left || dir == mathobjects.Right {
		kr.collisionHorizontal = dir
	} else {
		kr.collisionHorizontal = dir
		kr.collisionVertical = dir
	}
}

func (kr *KinematicRect) ApplyVelToPos(dtMs float64) {
	kr.Position.Add(kr.Velocity.Scale(dtMs))
}

func (kr KinematicRect) IsCollisionInDirection(dir mathobjects.Direction) bool {
	return kr.collisionHorizontal == dir || kr.collisionVertical == dir
}

func (kr KinematicRect) DoApplyForce(force ForceType) bool {
	return kr.ForcesToApply.Has(force)
}

func (kr KinematicRect) GetVelocity() mathobjects.Vector2[float64] {
	return kr.Velocity
}
