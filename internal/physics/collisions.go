package physics

import "math"

func collideKinematicWithStatic(kinematic *KinematicRect, static *Rect, dtMs float64) {
	if !static.DoApplyCollisions {
		return
	}
	kinematic.SetCollisionDir(None)

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

			kinematic.SetCollisionDir(Down)
		}
	} else if velY < 0 {
		// up
		futureKinematic.Size.Y += math.Abs(velY)
		futureKinematic.Position.Y -= math.Abs(velY)

		if futureKinematic.Intersects(*static) {
			kinematic.Position.Y = static.Position.Y + static.Size.Y
			kinematic.Velocity.Y = 0

			kinematic.SetCollisionDir(Up)
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

			kinematic.SetCollisionDir(Right)
		}
	} else if velX < 0 {
		// Left
		futureKinematic.Size.X += math.Abs(velX)
		futureKinematic.Position.X -= math.Abs(velX)

		if futureKinematic.Intersects(*static) {
			kinematic.Position.X = static.Position.X + static.Size.X
			kinematic.Velocity.X = 0

			kinematic.SetCollisionDir(Left)
		}
	}
}
