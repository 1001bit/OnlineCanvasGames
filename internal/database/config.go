package database

import "github.com/1001bit/OnlineCanvasGames/pkg/env"

type Config struct {
	user string
	name string
	pass string
	host string
}

func NewReadyConfig() *Config {
	return &Config{
		user: env.GetEnvVal("DB_USER"),
		name: env.GetEnvVal("DB_NAME"),
		pass: env.GetEnvVal("DB_PASS"),
		host: env.GetEnvVal("DB_HOST"),
	}
}
