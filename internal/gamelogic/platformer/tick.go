package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/gamelogic"
)

func (gl *PlatformerGL) tick(dtMs float64, writer gamelogic.RoomWriter) {
	gl.handleInput()
	deltas := gl.level.physEnv.Tick(dtMs, platformerConstants.Friction, platformerConstants.Gravity)
	if len(deltas) > 0 {
		writer.GlobalWriteMessage(gl.NewDeltasMessage(deltas))
	}
}
