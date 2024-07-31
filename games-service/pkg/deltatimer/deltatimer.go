package deltatimer

import "time"

type DeltaTimer struct {
	lastTick int64
}

func New() *DeltaTimer {
	return &DeltaTimer{
		lastTick: time.Now().UnixMicro(),
	}
}

func (timer *DeltaTimer) GetDeltaTimeMs() float64 {
	now := time.Now().UnixMicro()
	dt := now - timer.lastTick
	timer.lastTick = now

	return float64(dt) / 1000.0
}
