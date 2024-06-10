package physics

type Vector2f struct {
	X float64
	Y float64
}

func NewVector2(x, y float64) *Vector2f {
	return &Vector2f{
		X: x,
		Y: y,
	}
}

func (v2 *Vector2f) SetPosition(x, y float64) {
	v2.X = x
	v2.Y = y
}

func (v2 *Vector2f) Add(otherV2 *Vector2f) {
	v2.X += otherV2.X
	v2.Y += otherV2.Y
}

func (v2 *Vector2f) GetProduct(num float64) *Vector2f {
	return NewVector2(v2.X*num, v2.Y*num)
}
