package physics

func collideKinematicWithStatic(kinematic *KinematicRect, static *Rect, dtMs float64) {
	if !static.DoApplyCollisions {
		return
	}
	kinematic.collisionDir = CollisionDirection{None, None}

	// "path", that rect is going to pass
	futureKinematic := MakeRect(kinematic.Position.X, kinematic.Position.Y, kinematic.Size.X, kinematic.Size.Y, false)

	// Getting final value of velocity, which will be added to kinematicRect
	velX := kinematic.Velocity.X * dtMs
	velY := kinematic.Velocity.Y * dtMs

	// Vertical
	if velY > 0 {
		// down
		futureKinematic.Size.Y += velY
		if futureKinematic.Intersects(*static) {
			kinematic.Position.Y = static.Position.Y - kinematic.Size.Y
			kinematic.Velocity.Y = 0

			kinematic.collisionDir.Vertical = Down
		}
	} else if velY < 0 {
		// up
		futureKinematic.Position.Y += velY
		futureKinematic.Size.Y -= velY
		if futureKinematic.Intersects(*static) {
			kinematic.Position.Y = static.Position.Y + static.Size.Y
			kinematic.Velocity.Y = 0

			kinematic.collisionDir.Vertical = Up
		}
	}

	futureKinematic = MakeRect(kinematic.Position.X, kinematic.Position.Y, kinematic.Size.X, kinematic.Size.Y, false)

	// Horizontal
	if velX > 0 {
		// Right
		futureKinematic.Size.X += velX
		if futureKinematic.Intersects(*static) {
			kinematic.Position.X = static.Position.X - kinematic.Size.X
			kinematic.Velocity.X = 0

			kinematic.collisionDir.Horizontal = Right
		}
	} else if velX < 0 {
		// Left
		futureKinematic.Position.X += velX
		futureKinematic.Size.X -= velX
		if futureKinematic.Intersects(*static) {
			kinematic.Position.X = static.Position.X + static.Size.X
			kinematic.Velocity.X = 0

			kinematic.collisionDir.Horizontal = Left
		}
	}
}
