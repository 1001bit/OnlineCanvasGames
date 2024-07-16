package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/gamelogic"
)

func (gl *PlatformerGL) tick(dtMs float64, writer gamelogic.RoomWriter) {
	gl.handleInput()
	gl.level.physEng.Tick(dtMs, platformerConstants.Friction, platformerConstants.Gravity)
}
