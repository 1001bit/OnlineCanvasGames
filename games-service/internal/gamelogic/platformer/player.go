package platformer

import "github.com/1001bit/ocg-games-service/internal/mathobjects"

type Player struct {
	mathobjects.Rect

	velocity mathobjects.Vector2[float64]

	collisionVertical   mathobjects.Direction
	collisionHorizontal mathobjects.Direction

	futurePath mathobjects.Rect
}

func NewPlayer(rectID rectID) *Player {
	innerRect := mathobjects.CreateRect(100*float64(rectID), 100, 100, 100)

	return &Player{
		Rect: innerRect,

		velocity: mathobjects.Vector2[float64]{
			X: 0,
			Y: 0,
		},

		collisionVertical:   mathobjects.None,
		collisionHorizontal: mathobjects.None,

		futurePath: mathobjects.CreateRect(0, 0, 0, 0),
	}
}

func (p *Player) Control(speed, jump float64, input *PlayerInput) {
	if input == nil {
		return
	}

	if coeff, ok := input.GetHoldCoeff("left"); ok {
		p.velocity.X -= speed * coeff
	}

	if coeff, ok := input.GetHoldCoeff("right"); ok {
		p.velocity.X += speed * coeff
	}

	if input.IsHeld("jump") && p.collisionVertical == mathobjects.Down {
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
