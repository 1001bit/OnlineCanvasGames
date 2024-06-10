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
			gl.level.physEnv.Tick(float64(time.Since(lastTick)) / 1000000)

			// TODO: Send deltas instead of full level
			// TODO: Add teleport parameter to rect, if no interpolation is needed
			writer.GlobalWriteMessage(gl.NewFullLevelMessage())

			lastTick = time.Now()
		case <-doneChan:
			return
		}
	}
}
