package physics

type Environment struct {
	rects map[int]*Rect
}

func NewEnvironment() *Environment {
	return &Environment{
		rects: make(map[int]*Rect),
	}
}

func (e *Environment) InsertRect(r *Rect, id int) {
	e.rects[id] = r
}

func (e *Environment) GetRects() map[int]*Rect {
	return e.rects
}
