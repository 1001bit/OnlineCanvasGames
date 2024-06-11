package platformer

import (
	"time"

	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/gamelogic"
)

func (gl *PlatformerGL) handleInput() {
	for {
		select {
		case input := <-gl.inputChan:
			gl.level.ControlPlayerRect(input)
		default:
			return
		}
	}
}

func (gl *PlatformerGL) gameLoop(doneChan <-chan struct{}, writer gamelogic.RoomWriter) {
	ticker := time.NewTicker(time.Second / time.Duration(gl.tps))
	defer ticker.Stop()

	lastTick := time.Now()

	for {
		select {
		case <-ticker.C:
			gl.handleInput()
			dtMs := float64(time.Since(lastTick)) / 1000000

			gl.level.physEnv.Tick(dtMs)
			// TODO: send deltas instead of full level
			writer.GlobalWriteMessage(gl.NewLevelMessage())

			lastTick = time.Now()
		case <-doneChan:
			return
		}
	}
}
