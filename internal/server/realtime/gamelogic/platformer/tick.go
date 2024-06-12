package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/gamelogic"
)

func (gl *PlatformerGL) tick(dtMs float64, writer gamelogic.RoomWriter) {
	var (
		friction = 0.92
		gForce   = 1.0
	)

	gl.handleInput()
	deltas := gl.level.physEnv.Tick(dtMs, friction, gForce)
	writer.GlobalWriteMessage(gl.NewDeltasMessage(deltas))
}
