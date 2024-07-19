package physics

import "github.com/1001bit/OnlineCanvasGames/internal/mathobjects"

func (e *Engine) Tick(dtMs float64, constants PhysicsConstants) map[int]mathobjects.Vector2[float64] {
	positionsChanged := make(map[int]mathobjects.Vector2[float64])

	for id, kRect := range e.kinematicRects {
		startPos := kRect.GetPhysicalRect().Position

		applyGravityToVel(kRect, dtMs, constants.Gravity)
		applyFrictionToVel(kRect, constants.Friction)
		e.applyCollisions(kRect, dtMs)
		applyVelToPos(kRect, dtMs)

		// kRect.Velocity.RoundToZero(0.0001)

		if kRect.GetPhysicalRect().Position != startPos {
			positionsChanged[id] = kRect.GetPhysicalRect().Position
		}
	}

	return positionsChanged
}

func applyGravityToVel(rect *KinematicRect, dtMs, force float64) {
	if !rect.DoApplyGravity {
		return
	}

	rect.Velocity.Y += force * dtMs
}

func applyFrictionToVel(rect *KinematicRect, friction float64) {
	if !rect.DoApplyFriction {
		return
	}
	// for non gravitable rects
	if !rect.DoApplyGravity {
		rect.Velocity.Add(rect.Velocity.Scale(-friction))
		return
	}

	rect.Velocity.X -= rect.Velocity.X * friction
}

func (e *Engine) applyCollisions(kRect *KinematicRect, dtMs float64) {
	if !kRect.DoApplyCollisions {
		return
	}

	for _, rect := range e.staticRects {
		collideKinematicWithStatic(kRect, rect, dtMs)
	}
}

func applyVelToPos(rect *KinematicRect, dtMs float64) {
	rect.Position.Add(rect.Velocity.Scale(dtMs))
}
