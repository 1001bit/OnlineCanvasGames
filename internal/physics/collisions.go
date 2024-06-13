package physics

func collideKinematicWithSolid(kinematic *KinematicRect, solid *Rect, dtMs float64) {
	if !kinematic.Rect.doApplyCollisions || !solid.doApplyCollisions {
		return
	}
	kinematic.collisionDir = CollisionDirection{None, None}

	// "path", that rect is going to pass
	futureKinematic := MakeRect(kinematic.Position.X, kinematic.Position.Y, kinematic.Size.X, kinematic.Size.Y, false)

	// Getting final value of velocity, which will be added to kinematicRect
	velX := kinematic.Velocity.X * dtMs
	velY := kinematic.Velocity.Y * dtMs

	// Vertical
	if kinematic.Velocity.Y > 0 {
		// down
		futureKinematic.Size.Y += velY
		if futureKinematic.Intersects(*solid) {
			kinematic.Position.Y = solid.Position.Y - kinematic.Size.Y
			kinematic.Velocity.Y = 0

			kinematic.collisionDir.Vertical = Down
		}
	} else if kinematic.Velocity.Y < 0 {
		// up
		futureKinematic.Position.Y += velY
		futureKinematic.Size.Y -= kinematic.Velocity.Y
		if futureKinematic.Intersects(*solid) {
			kinematic.Position.Y = solid.Position.Y + solid.Size.Y
			kinematic.Velocity.Y = 0

			kinematic.collisionDir.Vertical = Up
		}
	}

	futureKinematic = MakeRect(kinematic.Position.X, kinematic.Position.Y, kinematic.Size.X, kinematic.Size.Y, false)

	// Horizontal
	if kinematic.Velocity.X > 0 {
		// Right
		futureKinematic.Size.X += velX
		if futureKinematic.Intersects(*solid) {
			kinematic.Position.X = solid.Position.X - kinematic.Size.X
			kinematic.Velocity.X = 0

			kinematic.collisionDir.Horizontal = Right
		}
	} else if kinematic.Velocity.X < 0 {
		// Left
		futureKinematic.Position.X += velX
		futureKinematic.Size.X -= kinematic.Velocity.X
		if futureKinematic.Intersects(*solid) {
			kinematic.Position.X = solid.Position.X + solid.Size.X
			kinematic.Velocity.X = 0

			kinematic.collisionDir.Horizontal = Left
		}
	}
}

func (e *Environment) applyCollisions(kRect *KinematicRect, dtMs float64) {
	for _, rect := range e.staticRects {
		collideKinematicWithSolid(kRect, rect, dtMs)
	}
}
