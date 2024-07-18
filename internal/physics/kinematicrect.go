package physics

// Collision direction
type Direction uint8

const (
	None  Direction = 0
	Left  Direction = 1
	Right Direction = 2
	Up    Direction = 3
	Down  Direction = 4
)

type CollisionDirection struct {
	Vertical   Direction
	Horizontal Direction
}

// Kinematic Rect
type KinematicRect struct {
	Rect

	Velocity Vector2f `json:"velocity"`

	collisionDir CollisionDirection

	DoApplyGravity  bool `json:"doApplyGravity"`
	DoApplyFriction bool `json:"doApplyFriction"`
}

func NewKinematicRect(rect Rect, doGravity, doFriction bool) *KinematicRect {
	return &KinematicRect{
		Rect: rect,

		Velocity: Vector2f{0, 0},

		collisionDir: CollisionDirection{
			Vertical:   None,
			Horizontal: None,
		},

		DoApplyGravity:  doGravity,
		DoApplyFriction: doFriction,
	}
}

func (kr *KinematicRect) AddToVel(add Vector2f) {
	kr.Velocity.Add(add)
}

func (kr *KinematicRect) SetCollisionDir(dir Direction) {
	if dir == Down || dir == Up {
		kr.collisionDir.Vertical = dir
	} else if dir == Left || dir == Right {
		kr.collisionDir.Horizontal = dir
	} else {
		kr.collisionDir.Horizontal = dir
		kr.collisionDir.Vertical = dir
	}
}

func (kr *KinematicRect) GetRect() Rect {
	return kr.Rect
}

func (kr *KinematicRect) GetCollisionDir() CollisionDirection {
	return kr.collisionDir
}
