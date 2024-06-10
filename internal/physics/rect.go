package physics

type Rect struct {
	Y float64 `json:"y"`
	X float64 `json:"x"`
	W float64 `json:"w"`
	H float64 `json:"h"`
}

func NewRect(x, y, w, h float64) *Rect {
	return &Rect{
		X: x,
		Y: y,
		W: w,
		H: h,
	}
}

type KinematicRect struct {
	rect *Rect

	velocity     *Vector2f
	acceleration *Vector2f

	doApplyGravity    bool
	doApplyCollisions bool
	doApplyFriction   bool
}

func NewKinematicRect(rect *Rect, gravity, collisions, friction bool) *KinematicRect {
	return &KinematicRect{
		rect:         rect,
		velocity:     NewVector2(0, 0),
		acceleration: NewVector2(0, 0),

		doApplyGravity:    gravity,
		doApplyCollisions: collisions,
		doApplyFriction:   friction,
	}
}

func (kr *KinematicRect) AddToAccel(add *Vector2f) {
	kr.acceleration.Add(add)
}

func (kr *KinematicRect) GetRect() *Rect {
	return kr.rect
}

func (kr *KinematicRect) applyGravityToAccel(dtMs, force float64) {
	if !kr.doApplyGravity {
		return
	}

	kr.acceleration.Y += force * dtMs
}

func (kr *KinematicRect) applyFrictionToVel(friction float64) {
	if !kr.doApplyFriction {
		return
	}

	if kr.doApplyGravity {
		// TODO: friction with collidables
		return
	}

	kr.velocity.X -= kr.velocity.X * friction
	kr.velocity.Y -= kr.velocity.Y * friction
}

func (kr *KinematicRect) applyAccelToVel(dtMs float64) {
	kr.velocity.Add(kr.acceleration.GetProduct(dtMs))
	kr.acceleration.X = 0
	kr.acceleration.Y = 0
}

func (kr *KinematicRect) applyVelToPos() {
	kr.rect.X += kr.velocity.X
	kr.rect.Y += kr.velocity.Y
}
