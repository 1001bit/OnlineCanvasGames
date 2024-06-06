package gamelogic

type Rect struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

func NewRect(x, y, w, h int) *Rect {
	return &Rect{
		X: x,
		Y: y,
		W: w,
		H: h,
	}
}
