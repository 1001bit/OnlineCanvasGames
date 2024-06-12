package physics

// TODO: collisions

func (e *Environment) Tick(dtMs, friction, gForce float64) map[int]*KinematicRect {
	deltas := make(map[int]*KinematicRect)

	for id, kRect := range e.kinematicRects {
		kRect.applyGravityToVel(dtMs, gForce)
		kRect.applyVelToPos(dtMs)
		kRect.applyFrictionToVel(friction)
		kRect.Velocity.RoundToZero(0.0001)

		if kRect.Velocity.X != 0 || kRect.Velocity.Y != 0 {
			deltas[id] = kRect
		}
	}

	return deltas
}
