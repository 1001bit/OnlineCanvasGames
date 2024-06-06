package gamelogic

type Level struct {
	rects       map[int]*Rect
	publicRects map[int]*Rect
}

func NewLevel() *Level {
	return &Level{
		rects:       make(map[int]*Rect),
		publicRects: make(map[int]*Rect),
	}
}

func (l *Level) InsertRect(r *Rect, id int, public bool) {
	l.rects[id] = r
	if public {
		l.publicRects[id] = r
	}
}

func (l *Level) GetAllRects() map[int]*Rect {
	return l.rects
}

func (l *Level) GetPublicRects() map[int]*Rect {
	return l.publicRects
}
