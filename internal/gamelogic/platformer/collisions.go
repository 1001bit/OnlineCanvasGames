package platformer

import (
	"math"

	"github.com/1001bit/OnlineCanvasGames/internal/mathobjects"
)

func (p *Player) DetectHorizontalCollision(block *Block, dtMs float64) mathobjects.Direction {
	if p.velocity.X == 0 {
		return mathobjects.None
	}

	// "path", that rect is going to pass
	futurePlayer := p.Rect
	finalVelX := p.velocity.X * dtMs

	if finalVelX > 0 {
		// right
		futurePlayer.Size.X += finalVelX

		if futurePlayer.Intersects(block.Rect) {
			return mathobjects.Right
		}
	} else {
		// left
		futurePlayer.Size.X += math.Abs(finalVelX)
		futurePlayer.Position.X -= math.Abs(finalVelX)

		if futurePlayer.Intersects(block.Rect) {
			return mathobjects.Left
		}
	}

	return mathobjects.None
}

func (p *Player) DetectVerticalCollision(block *Block, dtMs float64) mathobjects.Direction {
	if p.velocity.Y == 0 {
		return mathobjects.None
	}

	// "path", that rect is going to pass
	futurePlayer := p.Rect
	finalVelY := p.velocity.Y * dtMs

	if finalVelY > 0 {
		// down
		futurePlayer.Size.Y += finalVelY

		if futurePlayer.Intersects(block.Rect) {
			return mathobjects.Down
		}
	} else {
		// up
		futurePlayer.Size.Y += math.Abs(finalVelY)
		futurePlayer.Position.Y -= math.Abs(finalVelY)

		if futurePlayer.Intersects(block.Rect) {
			return mathobjects.Up
		}
	}

	return mathobjects.None
}

func (p *Player) ResolveCollision(block *Block, dir mathobjects.Direction) {
	if dir == mathobjects.None {
		return
	}

	p.SetCollisionDir(dir)

	switch dir {
	case mathobjects.Down:
		p.velocity.Y = 0
		p.Position.Y = block.GetPosition().Y - p.GetSize().Y
	case mathobjects.Up:
		p.velocity.Y = 0
		p.Position.Y = block.GetPosition().Y + block.GetSize().Y

	case mathobjects.Right:
		p.velocity.X = 0
		p.Position.X = block.GetPosition().X - p.GetSize().X
	case mathobjects.Left:
		p.velocity.X = 0
		p.Position.X = block.GetPosition().X + block.GetSize().X
	}
}
