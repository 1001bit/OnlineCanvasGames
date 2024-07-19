package refreshtoken

import "time"

const (
	// Duration of time when refresh token works
	ExpTime = time.Hour * 24 * 7 // 7d
	// Name of refresh token's cookie
	Name = "refresh"
)
