package physics

import "github.com/1001bit/OnlineCanvasGames/internal/mathobjects"

type PhysicalRect struct {
	mathobjects.Rect

	Collide bool `json:"canCollide"`
}

func MakePhysicalRect(x, y, w, h float64, canCollide bool) PhysicalRect {
	return PhysicalRect{
		Rect: mathobjects.MakeRect(x, y, w, h),

		Collide: canCollide,
	}
}

func NewPhysicalRect(x, y, w, h float64, canCollide bool) *PhysicalRect {
	rect := MakePhysicalRect(x, y, w, h, canCollide)
	return &rect
}

func (rect *PhysicalRect) SetPosition(x, y float64) {
	rect.Rect.Position.X = x
	rect.Rect.Position.Y = y
}

func (rect PhysicalRect) GetRect() mathobjects.Rect {
	return rect.Rect
}

func (rect PhysicalRect) CanCollide() bool {
	return rect.Collide
}
