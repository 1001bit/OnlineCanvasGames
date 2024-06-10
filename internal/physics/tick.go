package physics

func (e *Environment) Tick(dtMs float64) {
	for _, kRect := range e.kinematicRects {
		// TODO: Apply collision with other ticks

		kRect.applyGravityToAccel(dtMs, e.gForce)
		kRect.applyAccelToVel(dtMs)
		kRect.applyVelToPos()
		kRect.applyFrictionToVel(e.friction)
	}
}
