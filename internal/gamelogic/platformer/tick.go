package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/gamelogic"
)

func (gl *PlatformerGL) tick(dtMs float64, writer gamelogic.RoomWriter) {
	gl.handleInput()
	moved := gl.level.physEng.Tick(dtMs, platformerConstants.Physics)

	writer.GlobalWriteMessage(gl.NewUpdateMessage(moved))
}
