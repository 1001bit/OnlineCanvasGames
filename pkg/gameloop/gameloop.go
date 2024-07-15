package gameloop

import "time"

func Gameloop(callback func(dtMs float64), tps int, doneChan <-chan struct{}) {
	ticker := time.NewTicker(time.Second / time.Duration(tps))
	defer ticker.Stop()

	lastTick := time.Now()

	for {
		select {
		case <-ticker.C:
			dtMs := float64(time.Since(lastTick)) / 1000000
			lastTick = time.Now()
			callback(dtMs)

		case <-doneChan:
			return
		}
	}
}
