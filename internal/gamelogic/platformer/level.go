package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/gamelogic"
	"github.com/1001bit/OnlineCanvasGames/internal/mathobjects"
	"github.com/1001bit/OnlineCanvasGames/internal/physics"
)

type Level struct {
	// [rectID]rect
	playersRects map[int]*physics.KinematicRect
	levelRects   map[int]*physics.PhysicalRect

	// [userID]rectID
	userRectIDs map[int]int
}

func NewPlatformerLevel() *Level {
	level := &Level{
		playersRects: make(map[int]*physics.KinematicRect),
		levelRects:   make(map[int]*physics.PhysicalRect),

		userRectIDs: make(map[int]int),
	}

	block := physics.NewPhysicalRect(0, 500, 1000, 100, true)
	level.levelRects[10] = block

	return level
}

func (l *Level) Tick(dtMs float64, fullInputMap map[int]gamelogic.InputMap) map[int]mathobjects.Vector2[float64] {
	moved := make(map[int]mathobjects.Vector2[float64])

	// control player rects
	for userID, inputMap := range fullInputMap {
		l.controlPlayerRect(userID, inputMap)
	}

	// apply physics on player rects
	for rectID, kinematicRect := range l.playersRects {
		startPos := kinematicRect.GetPosition()

		// apply forces
		physics.ApplyForcesOn(kinematicRect, dtMs, globalPlatformerConstants.Physics)

		// apply collisions
		if kinematicRect.CanCollide() {
			kinematicRect.SetCollisionDir(mathobjects.None)
			for _, staticRect := range l.levelRects {
				physics.CollideKinematicWithStatic(kinematicRect, staticRect, dtMs)
			}
		}

		// move rect
		kinematicRect.ApplyVelToPos(dtMs)
		if startPos != kinematicRect.GetPosition() {
			moved[rectID] = kinematicRect.GetPosition()
		}
	}

	return moved
}

func (l *Level) controlPlayerRect(userID int, inputMap gamelogic.InputMap) {
	rectID := l.userRectIDs[userID]
	playerKRect, ok := l.playersRects[rectID]
	if !ok {
		return
	}

	addX := 0.0
	addY := 0.0

	if coeff, ok := inputMap.GetControlCoeff("left"); ok {
		addX -= globalPlatformerConstants.PlayerSpeed * coeff
	}

	if coeff, ok := inputMap.GetControlCoeff("right"); ok {
		addX += globalPlatformerConstants.PlayerSpeed * coeff
	}

	if inputMap.IsHeld("jump") && playerKRect.IsCollisionInDirection(mathobjects.Down) {
		addY -= globalPlatformerConstants.PlayerJump
	}

	playerKRect.AddToVel(addX, addY)
}
