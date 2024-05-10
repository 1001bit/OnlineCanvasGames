package basetoken

import "github.com/1001bit/OnlineCanvasGames/internal/env"

var secret = []byte("")

func InitJWTSecret() {
	secret = []byte(env.GetEnvVal("JWT_SECRET"))
}
