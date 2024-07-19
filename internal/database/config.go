package database

import "github.com/1001bit/OnlineCanvasGames/pkg/env"

type Config struct {
	user string
	name string
	pass string
}

func NewConfig() *Config {
	return &Config{
		env.GetEnvVal("DB_USER"),
		env.GetEnvVal("DB_NAME"),
		env.GetEnvVal("DB_PASS"),
	}
}
