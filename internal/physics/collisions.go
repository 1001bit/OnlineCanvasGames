package physics

import (
	"math"

	"github.com/1001bit/OnlineCanvasGames/internal/mathobjects"
)

func collideKinematicWithStatic(kinematic *KinematicRect, static *PhysicalRect, dtMs float64) {
	if !static.DoApplyCollisions {
		return
	}
	kinematic.SetCollisionDir(mathobjects.None)

	// "path", that rect is going to pass
	futureKinematic := kinematic.GetRect()

	// Getting final value of velocity, which will be added to kinematicRect
	finalVel := kinematic.GetVelocity().Scale(dtMs)

	// Vertical
	if finalVel.Y > 0 {
		// down
		futureKinematic.Size.Y += finalVel.Y

		if futureKinematic.Intersects(static.Rect) {
			kinematic.Position.Y = static.Position.Y - kinematic.Size.Y
			kinematic.Velocity.Y = 0

			kinematic.SetCollisionDir(mathobjects.Down)
		}
	} else if finalVel.Y < 0 {
		// up
		futureKinematic.Size.Y += math.Abs(finalVel.Y)
		futureKinematic.Position.Y -= math.Abs(finalVel.Y)

		if futureKinematic.Intersects(static.Rect) {
			kinematic.Position.Y = static.Position.Y + static.Size.Y
			kinematic.Velocity.Y = 0

			kinematic.SetCollisionDir(mathobjects.Up)
		}
	}

	futureKinematic = kinematic.GetRect()

	// Horizontal
	if finalVel.X > 0 {
		// Right
		futureKinematic.Size.X += finalVel.X

		if futureKinematic.Intersects(static.Rect) {
			kinematic.Position.X = static.Position.X - kinematic.Size.X
			kinematic.Velocity.X = 0

			kinematic.SetCollisionDir(mathobjects.Right)
		}
	} else if finalVel.X < 0 {
		// Left
		futureKinematic.Size.X += math.Abs(finalVel.X)
		futureKinematic.Position.X -= math.Abs(finalVel.X)

		if futureKinematic.Intersects(static.Rect) {
			kinematic.Position.X = static.Position.X + static.Size.X
			kinematic.Velocity.X = 0

			kinematic.SetCollisionDir(mathobjects.Left)
		}
	}
}
