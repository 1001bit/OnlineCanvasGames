package physics

type Rect struct {
	Position Vector2f `json:"position"`
	Size     Vector2f `json:"size"`
}

func MakeRect(x, y, w, h float64) Rect {
	return Rect{
		Position: Vector2f{x, y},
		Size:     Vector2f{w, h},
	}
}

type KinematicRect struct {
	Rect

	Velocity     Vector2f `json:"velocity"`
	acceleration Vector2f

	doApplyGravity    bool
	doApplyCollisions bool
	doApplyFriction   bool
}

func NewKinematicRect(rect Rect, gravity, collisions, friction bool) *KinematicRect {
	return &KinematicRect{
		Rect: rect,

		Velocity:     Vector2f{0, 0},
		acceleration: Vector2f{0, 0},

		doApplyGravity:    gravity,
		doApplyCollisions: collisions,
		doApplyFriction:   friction,
	}
}

func (kr *KinematicRect) AddToAccel(add Vector2f) {
	kr.acceleration.Add(add)
}

func (kr *KinematicRect) GetRect() Rect {
	return kr.Rect
}

func (kr *KinematicRect) applyGravityToAccel(dtMs, force float64) {
	if !kr.doApplyGravity {
		return
	}

	kr.acceleration.Y += force * dtMs
}

func (kr *KinematicRect) applyAccelToVel() {
	kr.Velocity.Add(kr.acceleration)
	kr.acceleration.X = 0
	kr.acceleration.Y = 0
}

func (kr *KinematicRect) applyVelToPos(dtMs float64) {
	kr.Rect.Position.Add(kr.Velocity.GetProduct(dtMs))
}
