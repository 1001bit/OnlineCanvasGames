package database

import (
	"database/sql"
	"fmt"

	"github.com/1001bit/overenv"
	_ "github.com/lib/pq"
)

type Config struct {
	User string `env:"DB_USER"`
	Name string `env:"DB_NAME"`
	Pass string `env:"DB_PASS"`
	Host string `env:"DB_HOST"`
	Port string `env:"DB_PORT"`
}

func NewFromEnv() (*sql.DB, error) {
	config := Config{}
	err := overenv.LoadStruct(&config)
	if err != nil {
		return nil, err
	}

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.Host, config.User, config.Pass, config.Name, config.Port)

	return sql.Open("postgres", connStr)
}
