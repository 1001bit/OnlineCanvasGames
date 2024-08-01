package platformer

import "github.com/1001bit/onlinecanvasgames/services/games/internal/mathobjects"

type Block struct {
	mathobjects.Rect
}

func NewBlock(x, y, w, h float64) *Block {
	return &Block{
		Rect: mathobjects.CreateRect(x, y, w, h),
	}
}
