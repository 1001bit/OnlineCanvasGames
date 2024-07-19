package physics

import "github.com/1001bit/OnlineCanvasGames/internal/mathobjects"

// Kinematic Rect
type KinematicRect struct {
	PhysicalRect

	Velocity mathobjects.Vector2[float64] `json:"velocity"`

	collisionVertical   mathobjects.Direction
	collisionHorizontal mathobjects.Direction

	DoApplyGravity  bool `json:"doApplyGravity"`
	DoApplyFriction bool `json:"doApplyFriction"`
}

func NewKinematicRect(rect PhysicalRect, doGravity, doFriction bool) *KinematicRect {
	return &KinematicRect{
		PhysicalRect: rect,

		Velocity: mathobjects.Vector2[float64]{
			X: 0,
			Y: 0,
		},

		collisionVertical:   mathobjects.None,
		collisionHorizontal: mathobjects.None,

		DoApplyGravity:  doGravity,
		DoApplyFriction: doFriction,
	}
}

func (kr *KinematicRect) AddToVel(add mathobjects.Vector2[float64]) {
	kr.Velocity.Add(add)
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

func (kr *KinematicRect) GetPhysicalRect() PhysicalRect {
	return kr.PhysicalRect
}

func (kr *KinematicRect) GetCollisionVertical() mathobjects.Direction {
	return kr.collisionVertical
}

func (kr *KinematicRect) GetCollisionHorizontal() mathobjects.Direction {
	return kr.collisionHorizontal
}
