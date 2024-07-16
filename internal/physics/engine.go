package physics

type Engine struct {
	staticRects    map[int]*Rect
	kinematicRects map[int]*KinematicRect
}

func NewEngine() *Engine {
	return &Engine{
		staticRects:    make(map[int]*Rect),
		kinematicRects: make(map[int]*KinematicRect),
	}
}

func (e *Engine) InsertRect(r *Rect, id int) {
	e.staticRects[id] = r
}

func (e *Engine) InsertKinematicRect(kr *KinematicRect, id int) {
	e.kinematicRects[id] = kr
}

func (e *Engine) DeleteRect(id int) {
	delete(e.staticRects, id)
	delete(e.kinematicRects, id)
}

func (e *Engine) GetStaticRects() map[int]*Rect {
	return e.staticRects
}

func (e *Engine) GetKinematicRects() map[int]*KinematicRect {
	return e.kinematicRects
}
