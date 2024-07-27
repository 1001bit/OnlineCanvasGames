package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/gamelogic"
	"github.com/1001bit/OnlineCanvasGames/internal/mathobjects"
	"github.com/1001bit/OnlineCanvasGames/pkg/concurrent"
	"github.com/1001bit/OnlineCanvasGames/pkg/fixedticker"
)

type rectID int

type LevelConfig struct {
	PlayerSpeed    float64 `json:"playerSpeed"`
	PlayerJump     float64 `json:"playerJump"`
	PlayerGravity  float64 `json:"playerGravity"`
	PlayerFriction float64 `json:"playerFriction"`
}

type Level struct {
	// [rectID]rect
	players map[rectID]*Player
	blocks  map[rectID]*Block

	// userID[data]
	playersData concurrent.ConcurrentMap[int, *PlayerData]

	config LevelConfig

	serverTPS float64
	clientTPS float64

	correctTicker *fixedticker.FixedTicker
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
		clientTPS = 40.0

		correctsPerSec = 1.0 / 5 // once per 5 sec
	)

	level := &Level{
		players: make(map[rectID]*Player),
		blocks:  make(map[rectID]*Block),

		playersData: concurrent.MakeMap[int, *PlayerData](),

		config: config,

		serverTPS: serverTPS,
		clientTPS: clientTPS,

		correctTicker: fixedticker.New(correctsPerSec),
	}

	block := NewBlock(0, 500, 1000, 100)
	level.blocks[10] = block

	return level
}

func (l *Level) Tick(dtMs float64, writer gamelogic.RoomWriter) {
	playersData, rUnlockFunc := l.playersData.GetMapForRead()
	defer rUnlockFunc()

	doCorrect := false
	l.correctTicker.Update(dtMs, func(_ float64) {
		doCorrect = true
	})

	// Physics
	// rectID[position]
	sentPlayers := make(map[rectID]mathobjects.Vector2[float64])

	for _, playerData := range playersData {
		player := playerData.player

		startPos := player.GetPosition()

		// Control
		playerData.ControlPlayer(l.config.PlayerSpeed, l.config.PlayerJump)

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

		// add player to movedPlayers if it moved or it's time for correction
		if startPos != player.GetPosition() || doCorrect {
			sentPlayers[playerData.rectID] = player.GetPosition()
		}
	}

	// Level Update Message
	writer.GlobalWriteMessage(NewLevelUpdateMessage(sentPlayers, doCorrect))
}
