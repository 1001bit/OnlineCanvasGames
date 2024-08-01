package database

import "github.com/1001bit/ocg-user-service/pkg/env"

type Config struct {
	user string
	name string
	pass string
	host string
	port string
}

func NewReadyConfig() *Config {
	return &Config{
		name: env.GetEnvVal("DB_NAME"),
		user: env.GetEnvVal("DB_USER"),
		pass: env.GetEnvVal("DB_PASS"),
		host: env.GetEnvVal("DB_HOST"),
		port: env.GetEnvVal("DB_PORT"),
	}
}
