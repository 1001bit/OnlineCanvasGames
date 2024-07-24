package fixedticker

type FixedTicker struct {
	accumulator float64
	tps         float64
}

func NewFixedTicker(tps float64) *FixedTicker {
	return &FixedTicker{
		accumulator: 0,
		tps:         tps,
	}
}

func (t *FixedTicker) Update(dt float64, callback func(fixedDtMs float64)) {
	t.accumulator += dt
	maxAccumulator := 1000 / t.tps
	for t.accumulator > maxAccumulator {
		callback(maxAccumulator)
		t.accumulator -= maxAccumulator
	}
}

func (t *FixedTicker) SetTPS(tps float64) {
	t.tps = tps
}

func (t *FixedTicker) GetTPS() float64 {
	return t.tps
}
