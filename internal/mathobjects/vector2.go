package mathobjects

type Vector2[T float64 | int] struct {
	X T `json:"x"`
	Y T `json:"y"`
}

func (v *Vector2[T]) SetPosition(x, y T) {
	v.X = x
	v.Y = y
}

func (v *Vector2[T]) Add(v2 Vector2[T]) {
	v.X += v2.X
	v.Y += v2.Y
}

func (v Vector2[T]) Scale(scalar T) Vector2[T] {
	return Vector2[T]{
		X: v.X * scalar,
		Y: v.Y * scalar,
	}
}
