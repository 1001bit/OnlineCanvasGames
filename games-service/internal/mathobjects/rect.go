package mathobjects

import "math"

type Rect struct {
	Position Vector2[float64] `json:"position"`
	Size     Vector2[float64] `json:"size"`
}

func CreateRect(x, y, w, h float64) Rect {
	return Rect{
		Position: Vector2[float64]{x, y},
		Size:     Vector2[float64]{w, h},
	}
}

func (rect *Rect) Intersects(rect2 Rect) bool {
	if rect.Position.X+rect.Size.X <= rect2.Position.X ||
		rect.Position.X >= rect2.Position.X+rect2.Size.X ||
		rect.Position.Y+rect.Size.Y <= rect2.Position.Y ||
		rect.Position.Y >= rect2.Position.Y+rect2.Size.Y {
		return false
	}

	return true
}

func (rect *Rect) SetPosition(x, y float64) {
	rect.Position.X = x
	rect.Position.Y = y
}

func (rect *Rect) SetSize(x, y float64) {
	rect.Size.X = x
	rect.Size.Y = y
}

func (rect *Rect) Extend(extX, extY float64) {
	rect.Size.X += math.Abs(extX)
	rect.Size.Y += math.Abs(extY)

	if extX < 0 {
		rect.Position.X -= math.Abs(extX)
	}
	if extY < 0 {
		rect.Position.Y -= math.Abs(extY)
	}
}

func (rect *Rect) GetPosition() Vector2[float64] {
	return rect.Position
}

func (rect *Rect) GetSize() Vector2[float64] {
	return rect.Size
}
