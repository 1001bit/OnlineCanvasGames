package physics

import "github.com/1001bit/OnlineCanvasGames/internal/mathobjects"

func (e *Engine) Tick(dtMs float64, constants PhysicsConstants) map[int]mathobjects.Vector2[float64] {
	positionsChanged := make(map[int]mathobjects.Vector2[float64])

	for id, kRect := range e.kinematicRects {
		startPos := kRect.GetPosition()

		applyGravityToVel(kRect, dtMs, constants.Gravity)
		applyFrictionToVel(kRect, constants.Friction)
		e.applyCollisions(kRect, dtMs)
		kRect.ApplyVelToPos(dtMs)

		if kRect.GetPosition() != startPos {
			positionsChanged[id] = kRect.GetPosition()
		}
	}

	return positionsChanged
}

func applyGravityToVel(rect *KinematicRect, dtMs, force float64) {
	if !rect.DoApplyGravity {
		return
	}

	rect.AddToVel(0, force*dtMs)
}

func applyFrictionToVel(rect *KinematicRect, friction float64) {
	if !rect.DoApplyFriction {
		return
	}

	rect.AddToVel(rect.GetVelocity().X*-friction, 0)
	// for non gravitable rects
	if !rect.DoApplyGravity {
		rect.AddToVel(0, rect.GetVelocity().Y*-friction)
	}
}

func (e *Engine) applyCollisions(kRect *KinematicRect, dtMs float64) {
	if !kRect.DoApplyCollisions {
		return
	}

	for _, rect := range e.staticRects {
		collideKinematicWithStatic(kRect, rect, dtMs)
	}
}
