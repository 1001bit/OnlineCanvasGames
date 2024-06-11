package physics

func (e *Environment) Tick(dtMs float64) {
	for _, kRect := range e.kinematicRects {
		kRect.applyGravityToAccel(dtMs, e.gForce)
		kRect.applyAccelToVel()
		kRect.Velocity.RoundToZero(0.0001)
		kRect.applyVelToPos(dtMs)

		// TODO: Friction and collisions

		// TODO: This + return all the rects, where velocity != 0
		// kRect.applyVelToPos(dtMs)
		// kRect.applyGravityToAccel(dtMs, e.gForce)
		// kRect.applyAccelToVel()
		// kRect.Velocity.RoundToZero(0.0001)
	}
}
