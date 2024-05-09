package basetoken

import "github.com/1001bit/OnlineCanvasGames/internal/env"

var secret = []byte("")

func InitJWTSecret() {
	secret = []byte(env.GetEnvVal("ACC_SECRET"))
}
