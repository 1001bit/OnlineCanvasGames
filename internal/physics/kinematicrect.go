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

	doApplyGravity  bool
	doApplyFriction bool
}

func NewKinematicRect(rect Rect, doGravity, doFriction bool) *KinematicRect {
	return &KinematicRect{
		Rect: rect,

		Velocity: Vector2f{0, 0},

		collisionDir: CollisionDirection{
			Vertical:   None,
			Horizontal: None,
		},

		doApplyGravity:  doGravity,
		doApplyFriction: doFriction,
	}
}

func (kr *KinematicRect) AddToVel(add Vector2f) {
	kr.Velocity.Add(add)
}

func (kr *KinematicRect) GetRect() Rect {
	return kr.Rect
}

func (kr *KinematicRect) GetCollisionDir() CollisionDirection {
	return kr.collisionDir
}
