package accesstoken

import "time"

const (
	// Time when token is valid
	ExpTime = time.Hour * 6 // 6h
	// Cookie name
	Name = "access"
)
