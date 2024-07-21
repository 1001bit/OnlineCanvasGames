package physics

import "github.com/1001bit/OnlineCanvasGames/internal/mathobjects"

type PhysicalRect struct {
	mathobjects.Rect

	DoApplyCollisions bool `json:"doApplyCollisions"`
}

func MakePhysicalRect(x, y, w, h float64, doApplyCollisions bool) PhysicalRect {
	return PhysicalRect{
		Rect: mathobjects.MakeRect(x, y, w, h),

		DoApplyCollisions: doApplyCollisions,
	}
}

func NewPhysicalRect(x, y, w, h float64, doApplyCollisions bool) *PhysicalRect {
	rect := MakePhysicalRect(x, y, w, h, doApplyCollisions)
	return &rect
}

func (rect PhysicalRect) GetRect() mathobjects.Rect {
	return rect.Rect
}
