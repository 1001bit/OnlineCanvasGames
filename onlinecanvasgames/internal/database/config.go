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
		user: env.GetEnvVal("POSTGRES_USER"),
		name: env.GetEnvVal("POSTGRES_DB"),
		pass: env.GetEnvVal("POSTGRES_PASSWORD"),
		host: env.GetEnvVal("DB_HOST"),
	}
}
