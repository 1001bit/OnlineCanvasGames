package platformer

import (
	"time"
)

func (gl *PlatformerGL) gameLoop(doneChan <-chan struct{}) {
	ticker := time.NewTicker(time.Second / time.Duration(gl.ticksPerSecond))
	defer ticker.Stop()

	lastTick := time.Now()

	for {
		select {
		case <-ticker.C:
			gl.level.physEnv.Tick(float64(time.Since(lastTick)) / 1000000)
			lastTick = time.Now()
		case <-doneChan:
			return
		}
	}
}
