package physics

type Rect struct {
	Position Vector2f `json:"position"`
	Size     Vector2f `json:"size"`

	DoApplyCollisions bool `json:"doApplyCollisions"`
}

func MakeRect(x, y, w, h float64, doApplyCollisions bool) Rect {
	return Rect{
		Position: Vector2f{x, y},
		Size:     Vector2f{w, h},

		DoApplyCollisions: doApplyCollisions,
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
