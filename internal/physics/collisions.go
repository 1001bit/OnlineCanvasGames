package physics

import (
	"math"

	"github.com/1001bit/OnlineCanvasGames/internal/mathobjects"
)

func CollideKinematicWithStatic(kinematic Kinematic, static Physical, dtMs float64) {
	if !static.CanCollide() {
		return
	}

	// TODO: skip if kinematic->solid direction == kinematic.collisionDir

	// "path", that rect is going to pass
	futureKinematic := kinematic.GetRect()

	// Getting final value of velocity, which will be added to kinematicRect
	finalVel := kinematic.GetVelocity().Scale(dtMs)

	// Vertical
	if finalVel.Y > 0 {
		// down
		futureKinematic.Size.Y += finalVel.Y

		if futureKinematic.Intersects(static.GetRect()) {
			newPosY := static.GetPosition().Y - kinematic.GetSize().Y

			kinematic.SetPosition(kinematic.GetPosition().X, newPosY)
			kinematic.SetVelocity(kinematic.GetVelocity().X, 0)

			kinematic.SetCollisionDir(mathobjects.Down)
		}
	} else if finalVel.Y < 0 {
		// up
		futureKinematic.Size.Y += math.Abs(finalVel.Y)
		futureKinematic.Position.Y -= math.Abs(finalVel.Y)

		if futureKinematic.Intersects(static.GetRect()) {
			newPosY := static.GetPosition().Y + static.GetSize().Y

			kinematic.SetPosition(kinematic.GetPosition().X, newPosY)
			kinematic.SetVelocity(kinematic.GetVelocity().X, 0)

			kinematic.SetCollisionDir(mathobjects.Up)
		}
	}

	futureKinematic = kinematic.GetRect()

	// Horizontal
	if finalVel.X > 0 {
		// Right
		futureKinematic.Size.X += finalVel.X

		if futureKinematic.Intersects(static.GetRect()) {
			newPosX := static.GetPosition().X - kinematic.GetSize().X

			kinematic.SetPosition(newPosX, kinematic.GetPosition().Y)
			kinematic.SetVelocity(0, kinematic.GetVelocity().Y)

			kinematic.SetCollisionDir(mathobjects.Right)
		}
	} else if finalVel.X < 0 {
		// Left
		futureKinematic.Size.X += math.Abs(finalVel.X)
		futureKinematic.Position.X -= math.Abs(finalVel.X)

		if futureKinematic.Intersects(static.GetRect()) {
			newPosX := static.GetPosition().X + static.GetSize().X

			kinematic.SetPosition(newPosX, kinematic.GetPosition().Y)
			kinematic.SetVelocity(0, kinematic.GetVelocity().Y)

			kinematic.SetCollisionDir(mathobjects.Left)
		}
	}
}
