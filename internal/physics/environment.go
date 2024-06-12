package physics

type Environment struct {
	staticRects    map[int]*Rect
	kinematicRects map[int]*KinematicRect
}

func NewEnvironment() *Environment {
	return &Environment{
		staticRects:    make(map[int]*Rect),
		kinematicRects: make(map[int]*KinematicRect),
	}
}

func (e *Environment) InsertRect(r *Rect, id int) {
	e.staticRects[id] = r
}

func (e *Environment) InsertKinematicRect(kr *KinematicRect, id int) {
	e.kinematicRects[id] = kr
}

func (e *Environment) DeleteRect(id int) {
	delete(e.staticRects, id)
	delete(e.kinematicRects, id)
}

func (e *Environment) GetStaticRects() map[int]*Rect {
	return e.staticRects
}

func (e *Environment) GetKinematicRects() map[int]*KinematicRect {
	return e.kinematicRects
}
