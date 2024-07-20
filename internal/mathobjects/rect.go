package mathobjects

type Rect struct {
	Position Vector2[float64] `json:"position"`
	Size     Vector2[float64] `json:"size"`
}

func MakeRect(x, y, w, h float64) Rect {
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

func (rect *Rect) GetPosition() Vector2[float64] {
	return rect.Position
}

func (rect *Rect) GetSize() Vector2[float64] {
	return rect.Size
}
