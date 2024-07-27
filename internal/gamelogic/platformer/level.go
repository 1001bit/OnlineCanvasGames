package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/gamelogic"
	"github.com/1001bit/OnlineCanvasGames/internal/mathobjects"
	"github.com/1001bit/OnlineCanvasGames/pkg/concurrent"
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

	// userID[data]
	playersData concurrent.ConcurrentMap[int, *PlayerData]

	config LevelConfig

	serverTPS float64
	clientTPS float64
}

func NewPlatformerLevel() *Level {
	var (
		config = LevelConfig{
			PlayerSpeed:    3,
			PlayerJump:     5,
			PlayerGravity:  0.03,
			PlayerFriction: 0.3,
		}

		serverTPS = 20.0
		clientTPS = 50.0
	)

	level := &Level{
		players: make(map[int]*Player),
		blocks:  make(map[int]*Block),

		playersData: concurrent.MakeMap[int, *PlayerData](),

		config: config,

		serverTPS: serverTPS,
		clientTPS: clientTPS,
	}

	block := NewBlock(0, 500, 1000, 100)
	level.blocks[10] = block

	return level
}

func (l *Level) Tick(dtMs float64, writer gamelogic.RoomWriter) {
	movedPlayers := make(map[int]mathobjects.Vector2[float64])

	// Controls
	playersData, rUnlockFunc := l.playersData.GetMapForRead()
	for _, playerData := range playersData {
		playerData.ControlPlayer(l.config.PlayerSpeed, l.config.PlayerJump)
	}
	rUnlockFunc()

	// Physics
	for rectID, player := range l.players {
		startPos := player.GetPosition()

		// Forces
		player.ApplyGravity(l.config.PlayerGravity, dtMs)
		player.ApplyFriction(l.config.PlayerFriction)

		// Collisions and displacement
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

		if startPos != player.GetPosition() {
			movedPlayers[rectID] = player.GetPosition()
		}
	}

	// Level Update Message
	writer.GlobalWriteMessage(NewLevelUpdateMessage(movedPlayers))
}
