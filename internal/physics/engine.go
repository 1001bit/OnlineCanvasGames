package physics

type Engine struct {
	staticRects    map[int]*PhysicalRect
	kinematicRects map[int]*KinematicRect
}

func NewEngine() *Engine {
	return &Engine{
		staticRects:    make(map[int]*PhysicalRect),
		kinematicRects: make(map[int]*KinematicRect),
	}
}

func (e *Engine) InsertStaticRect(r *PhysicalRect, id int) {
	e.staticRects[id] = r
}

func (e *Engine) InsertKinematicRect(kr *KinematicRect, id int) {
	e.kinematicRects[id] = kr
}

func (e *Engine) DeleteRect(id int) {
	delete(e.staticRects, id)
	delete(e.kinematicRects, id)
}

func (e *Engine) GetStaticRects() map[int]*PhysicalRect {
	return e.staticRects
}

func (e *Engine) GetKinematicRects() map[int]*KinematicRect {
	return e.kinematicRects
}
