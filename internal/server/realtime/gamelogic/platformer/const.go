package platformer

type Constants struct {
	Friction    float64 `json:"friction"`
	Gravity     float64 `json:"gravity"`
	PlayerSpeed float64 `json:"playerSpeed"`
	PlayerJump  float64 `json:"playerJump"`
}

var platformerConstants = Constants{
	Friction:    0.92,
	Gravity:     0.02,
	PlayerSpeed: 5,
	PlayerJump:  3,
}
