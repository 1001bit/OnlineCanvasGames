package physics

type Environment struct {
	rects          map[int]*Rect
	kinematicRects map[int]*KinematicRect

	friction float64
	gForce   float64
}

func NewEnvironment(friction, gForce float64) *Environment {
	return &Environment{
		rects:          make(map[int]*Rect),
		kinematicRects: make(map[int]*KinematicRect),

		friction: friction,
		gForce:   gForce,
	}
}

func (e *Environment) InsertRect(r *Rect, id int) {
	e.rects[id] = r
}

func (e *Environment) InsertKinematicRect(kr *KinematicRect, id int) {
	e.kinematicRects[id] = kr
}

func (e *Environment) DeleteRect(id int) {
	delete(e.rects, id)
	delete(e.kinematicRects, id)
}

func (e *Environment) GetRects() map[int]*Rect {
	return e.rects
}

func (e *Environment) GetKinematicRect(rectID int) (*KinematicRect, bool) {
	kr, ok := e.kinematicRects[rectID]
	return kr, ok
}
