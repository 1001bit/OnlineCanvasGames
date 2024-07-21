package physics

type ForceType string

const (
	GravityType  ForceType = "gravity"
	FrictionType ForceType = "friction"
)

func ApplyForcesOn(rect Kinematic, dtMs float64, constants PhysicsConstants) {
	if rect.DoApplyForce(GravityType) {
		applyGravityToVel(rect, dtMs, constants.Gravity)
	}
	if rect.DoApplyForce(FrictionType) {
		applyFrictionToVel(rect, constants.Friction)
	}
}

func applyGravityToVel(rect Kinematic, dtMs, force float64) {
	rect.AddToVel(0, force*dtMs)
}

func applyFrictionToVel(rect Kinematic, friction float64) {
	rect.AddToVel(rect.GetVelocity().X*-friction, 0)
	// for non gravitable rects
	if !rect.DoApplyForce(GravityType) {
		rect.AddToVel(0, rect.GetVelocity().Y*-friction)
	}
}
