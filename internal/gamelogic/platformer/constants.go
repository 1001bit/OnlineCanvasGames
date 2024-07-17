package platformer

import "github.com/1001bit/OnlineCanvasGames/internal/physics"

type Constants struct {
	Physics     physics.PhysicsConstants `json:"physics"`
	PlayerSpeed float64                  `json:"playerSpeed"`
	PlayerJump  float64                  `json:"playerJump"`
}

var platformerConstants = Constants{
	Physics: physics.PhysicsConstants{
		Friction: 0.92,
		Gravity:  0.02,
	},
	PlayerSpeed: 5,
	PlayerJump:  3,
}
