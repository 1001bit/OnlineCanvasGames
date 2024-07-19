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
