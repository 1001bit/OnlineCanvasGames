package platformer

import (
	"math"

	"github.com/1001bit/OnlineCanvasGames/internal/mathobjects"
)

func CollidePlayerWithBlock(player *Player, block *Block, dtMs float64) {
	// TODO: skip if player->block direction == player.collisionDir
	// TODO: use interfaces for collided objects

	// "path", that rect is going to pass
	futurePlayer := player.Rect

	// Getting final value of velocity, which will be added to kinematicRect
	finalVel := player.velocity.Scale(dtMs)

	// Vertical
	if finalVel.Y > 0 {
		// down
		futurePlayer.Size.Y += finalVel.Y

		if futurePlayer.Intersects(block.Rect) {
			player.Position.Y = block.GetPosition().Y - player.GetSize().Y
			player.velocity.Y = 0

			player.SetCollisionDir(mathobjects.Down)
		}
	} else if finalVel.Y < 0 {
		// up
		futurePlayer.Size.Y += math.Abs(finalVel.Y)
		futurePlayer.Position.Y -= math.Abs(finalVel.Y)

		if futurePlayer.Intersects(block.Rect) {
			player.Position.Y = block.GetPosition().Y + block.GetSize().Y
			player.velocity.Y = 0

			player.SetCollisionDir(mathobjects.Up)
		}
	}

	futurePlayer = player.Rect

	// Horizontal
	if finalVel.X > 0 {
		// Right
		futurePlayer.Size.X += finalVel.X

		if futurePlayer.Intersects(block.Rect) {
			player.Position.X = block.GetPosition().X - player.GetSize().X
			player.velocity.X = 0

			player.SetCollisionDir(mathobjects.Right)
		}
	} else if finalVel.X < 0 {
		// Left
		futurePlayer.Size.X += math.Abs(finalVel.X)
		futurePlayer.Position.X -= math.Abs(finalVel.X)

		if futurePlayer.Intersects(block.Rect) {
			player.Position.X = block.GetPosition().X + block.GetSize().X
			player.velocity.X = 0

			player.SetCollisionDir(mathobjects.Left)
		}
	}
}
