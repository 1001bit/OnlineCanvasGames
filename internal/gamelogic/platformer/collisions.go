package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/mathobjects"
)

func (p *Player) DetectHorizontalCollision(block *Block, dtMs float64) mathobjects.Direction {
	if p.velocity.X == 0 {
		return mathobjects.None
	}

	playerPath := p.Rect
	playerPath.Extend(p.velocity.X*dtMs, 0)

	if !playerPath.Intersects(block.Rect) {
		return mathobjects.None
	}

	if p.velocity.X > 0 {
		return mathobjects.Right
	} else {
		return mathobjects.Left
	}
}

func (p *Player) DetectVerticalCollision(block *Block, dtMs float64) mathobjects.Direction {
	if p.velocity.Y == 0 {
		return mathobjects.None
	}

	playerPath := p.Rect
	playerPath.Extend(0, p.velocity.Y*dtMs)

	if !playerPath.Intersects(block.Rect) {
		return mathobjects.None
	}

	if p.velocity.Y > 0 {
		return mathobjects.Down
	} else {
		return mathobjects.Up
	}
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
