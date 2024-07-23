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

	// Control
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

	// Physics
	for rectID, player := range l.players {
		startPos := player.GetPosition()

		// Forces
		player.ApplyGravity(l.config.PlayerGravity, dtMs)
		player.ApplyFriction(l.config.PlayerFriction)

		// Collisions and movement
		player.SetCollisionDir(mathobjects.None)

		// Horizontal
		for _, block := range l.blocks {
			dir := player.DetectHorizontalCollision(block, dtMs)
			if dir != mathobjects.None {
				player.ResolveCollision(block, dir)
				break
			}
		}
		player.Position.X += player.velocity.X * dtMs

		// Vertical
		for _, block := range l.blocks {
			dir := player.DetectVerticalCollision(block, dtMs)
			if dir != mathobjects.None {
				player.ResolveCollision(block, dir)
				break
			}
		}
		player.Position.Y += player.velocity.Y * dtMs

		// Register to moved rects
		if startPos != player.GetPosition() {
			moved[rectID] = player.GetPosition()
		}
	}

	return moved
}
