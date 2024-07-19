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
	futureKinematic := mathobjects.MakeRect(kinematic.Position.X, kinematic.Position.Y, kinematic.Size.X, kinematic.Size.Y)

	// Getting final value of velocity, which will be added to kinematicRect
	velX := kinematic.Velocity.X * dtMs
	velY := kinematic.Velocity.Y * dtMs

	// Vertical
	if velY > 0 {
		// down
		futureKinematic.Size.Y += velY

		if futureKinematic.Intersects(static.Rect) {
			kinematic.Position.Y = static.Position.Y - kinematic.Size.Y
			kinematic.Velocity.Y = 0

			kinematic.SetCollisionDir(mathobjects.Down)
		}
	} else if velY < 0 {
		// up
		futureKinematic.Size.Y += math.Abs(velY)
		futureKinematic.Position.Y -= math.Abs(velY)

		if futureKinematic.Intersects(static.Rect) {
			kinematic.Position.Y = static.Position.Y + static.Size.Y
			kinematic.Velocity.Y = 0

			kinematic.SetCollisionDir(mathobjects.Up)
		}
	}

	futureKinematic = mathobjects.MakeRect(kinematic.Position.X, kinematic.Position.Y, kinematic.Size.X, kinematic.Size.Y)

	// Horizontal
	if velX > 0 {
		// Right
		futureKinematic.Size.X += velX

		if futureKinematic.Intersects(static.Rect) {
			kinematic.Position.X = static.Position.X - kinematic.Size.X
			kinematic.Velocity.X = 0

			kinematic.SetCollisionDir(mathobjects.Right)
		}
	} else if velX < 0 {
		// Left
		futureKinematic.Size.X += math.Abs(velX)
		futureKinematic.Position.X -= math.Abs(velX)

		if futureKinematic.Intersects(static.Rect) {
			kinematic.Position.X = static.Position.X + static.Size.X
			kinematic.Velocity.X = 0

			kinematic.SetCollisionDir(mathobjects.Left)
		}
	}
}
