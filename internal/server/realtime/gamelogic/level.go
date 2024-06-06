package gamelogic

type Level struct {
	rects map[int]*Rect
}

func NewLevel(rects map[int]*Rect) *Level {
	return &Level{
		rects: rects,
	}
}

func (l *Level) GetRects(onlyDynamic, onlyPublic bool) map[int]*Rect {
	if !onlyDynamic && !onlyPublic {
		return l.rects
	}

	result := make(map[int]*Rect)

	for i := range l.rects {
		if onlyPublic && !l.rects[i].public {
			continue
		}
		if onlyDynamic && !l.rects[i].Dynamic {
			continue
		}
		result[i] = l.rects[i]
	}

	return result
}
