package physics

import "math"

type Vector2f struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (v2 *Vector2f) SetPosition(x, y float64) {
	v2.X = x
	v2.Y = y
}

func (v2 *Vector2f) Add(otherV2 Vector2f) {
	v2.X += otherV2.X
	v2.Y += otherV2.Y
}

func (v2 *Vector2f) RoundToZero(num float64) {
	if math.Abs(v2.X) <= num && v2.X != 0 {
		v2.X = 0
	}
	if math.Abs(v2.Y) <= num && v2.Y != 0 {
		v2.Y = 0
	}
}

func (v2 Vector2f) Scale(scalar float64) Vector2f {
	return Vector2f{
		X: v2.X * scalar,
		Y: v2.Y * scalar,
	}
}