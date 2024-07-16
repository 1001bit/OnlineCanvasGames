package physics

func (e *Engine) Tick(dtMs, friction, gForce float64) {
	for _, kRect := range e.kinematicRects {
		applyGravityToVel(kRect, dtMs, gForce)
		applyFrictionToVel(kRect, friction)
		e.applyCollisions(kRect, dtMs)
		applyVelToPos(kRect, dtMs)

		kRect.Velocity.RoundToZero(0.0001)
	}
}

func applyGravityToVel(rect *KinematicRect, dtMs, force float64) {
	if !rect.doApplyGravity {
		return
	}

	rect.Velocity.Y += force * dtMs
}

func applyFrictionToVel(rect *KinematicRect, friction float64) {
	if !rect.doApplyFriction {
		return
	}
	// for non gravitable rects
	if !rect.doApplyGravity {
		rect.Velocity.Add(rect.Velocity.Scale(-friction))
		return
	}

	rect.Velocity.X -= rect.Velocity.X * friction
}

func (e *Engine) applyCollisions(kRect *KinematicRect, dtMs float64) {
	for _, rect := range e.staticRects {
		collideKinematicWithSolid(kRect, rect, dtMs)
	}
}

func applyVelToPos(rect *KinematicRect, dtMs float64) {
	rect.Rect.Position.Add(rect.Velocity.Scale(dtMs))
}
