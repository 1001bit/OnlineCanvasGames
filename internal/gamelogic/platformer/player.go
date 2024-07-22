package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/gamelogic"
	"github.com/1001bit/OnlineCanvasGames/internal/mathobjects"
)

type Player struct {
	mathobjects.Rect

	velocity mathobjects.Vector2[float64]

	collisionVertical   mathobjects.Direction
	collisionHorizontal mathobjects.Direction
}

func NewPlayer(rectID int) *Player {
	innerRect := mathobjects.MakeRect(100*float64(rectID), 100, 100, 100)

	return &Player{
		Rect: innerRect,

		velocity: mathobjects.Vector2[float64]{
			X: 0,
			Y: 0,
		},

		collisionVertical:   mathobjects.None,
		collisionHorizontal: mathobjects.None,
	}
}

func (p *Player) Control(speed, jump float64, inputMap gamelogic.InputMap) {
	if coeff, ok := inputMap.GetControlCoeff("left"); ok {
		p.velocity.X -= speed * coeff
	}

	if coeff, ok := inputMap.GetControlCoeff("right"); ok {
		p.velocity.X += speed * coeff
	}

	if inputMap.IsHeld("jump") && p.collisionVertical == mathobjects.Down {
		p.velocity.Y -= jump
	}
}

func (p *Player) ApplyGravity(force, dtMs float64) {
	p.velocity.Y += force * dtMs
}

func (p *Player) ApplyFriction(force float64) {
	p.velocity.X *= force
	// p.velocity.Y *= force
}

func (p *Player) SetCollisionDir(dir mathobjects.Direction) {
	if dir == mathobjects.Down || dir == mathobjects.Up {
		p.collisionVertical = dir
	} else if dir == mathobjects.Left || dir == mathobjects.Right {
		p.collisionHorizontal = dir
	} else {
		p.collisionHorizontal = dir
		p.collisionVertical = dir
	}
}

func (p *Player) ApplyVelToPos(dtMs float64) {
	p.Position.Add(p.velocity.Scale(dtMs))
}
