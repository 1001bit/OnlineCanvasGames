package physics

func (e *Environment) Tick(dtMs, friction, gForce float64) map[int]*KinematicRect {
	deltas := make(map[int]*KinematicRect)

	for id, kRect := range e.kinematicRects {
		kRect.applyGravityToVel(dtMs, gForce)
		kRect.applyFrictionToVel(friction)

		kRect.applyVelToPos(dtMs)

		e.applyCollisions(kRect, dtMs)

		didRound := kRect.Velocity.RoundToZero(0.0001)

		if kRect.Velocity.X != 0 || kRect.Velocity.Y != 0 || didRound {
			deltas[id] = kRect
		}
	}

	return deltas
}
