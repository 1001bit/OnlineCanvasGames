package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/gamelogic"
	"github.com/1001bit/OnlineCanvasGames/internal/mathobjects"
)

type LevelConfig struct {
	PlayerSpeed    float64 `json:"playerSpeed"`
	PlayerJump     float64 `json:"playerJump"`
	PlayerGravity  float64 `json:"playerGravity"`
	PlayerFriction float64 `json:"playerFriction"`
}

type Level struct {
	// [rectID]rect
	players map[int]*Player
	blocks  map[int]*Block

	config LevelConfig

	// [userID]rectID
	userRectIDs map[int]int
}

func NewPlatformerLevel() *Level {
	config := LevelConfig{
		PlayerSpeed:    3,
		PlayerJump:     5,
		PlayerGravity:  0.03,
		PlayerFriction: 0.3,
	}

	level := &Level{
		players: make(map[int]*Player),
		blocks:  make(map[int]*Block),

		config: config,

		userRectIDs: make(map[int]int),
	}

	block := NewBlock(0, 500, 1000, 100)
	level.blocks[10] = block

	return level
}

func (l *Level) Tick(dtMs float64, fullInputMap map[int]gamelogic.InputMap) map[int]mathobjects.Vector2[float64] {
	moved := make(map[int]mathobjects.Vector2[float64])

	// TODO: Fixed Timestep

	// control player rects
	for userID, inputMap := range fullInputMap {
		rectID, ok := l.userRectIDs[userID]
		if !ok {
			continue
		}
		player, ok := l.players[rectID]
		if !ok {
			continue
		}

		player.Control(l.config.PlayerSpeed, l.config.PlayerJump, inputMap)
	}

	// apply physics on player rects
	for rectID, player := range l.players {
		startPos := player.GetPosition()

		// apply forces
		player.ApplyGravity(l.config.PlayerGravity, dtMs)
		player.ApplyFriction(l.config.PlayerFriction)

		// apply collision
		player.SetCollisionDir(mathobjects.None)
		for _, block := range l.blocks {
			CollidePlayerWithBlock(player, block, dtMs)
		}

		// move rect
		player.ApplyVelToPos(dtMs)
		if startPos != player.GetPosition() {
			moved[rectID] = player.GetPosition()
		}
	}

	return moved
}
