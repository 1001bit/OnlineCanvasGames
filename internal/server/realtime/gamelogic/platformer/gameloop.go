package platformer

import (
	"time"
)

func (gl *PlatformerGL) tick(dtMs float64) {

}

func (gl *PlatformerGL) gameLoop(doneChan <-chan struct{}) {
	ticker := time.NewTicker(time.Second / time.Duration(gl.ticksPerSecond))
	defer ticker.Stop()

	lastTick := time.Now()

	for {
		select {
		case <-ticker.C:
			gl.tick(float64(time.Since(lastTick)) / 1000000)
			lastTick = time.Now()
		case <-doneChan:
			return
		}
	}
}
