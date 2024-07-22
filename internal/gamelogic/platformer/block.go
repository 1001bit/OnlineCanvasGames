package platformer

import "github.com/1001bit/OnlineCanvasGames/internal/mathobjects"

type Block struct {
	mathobjects.Rect
}

func NewBlock(x, y, w, h float64) *Block {
	return &Block{
		Rect: mathobjects.MakeRect(x, y, w, h),
	}
}
