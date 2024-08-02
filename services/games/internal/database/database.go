package database

import (
	"database/sql"
	"fmt"

	"github.com/1001bit/overenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

type Config struct {
	User string `env:"DB_USER"`
	Name string `env:"DB_NAME"`
	Pass string `env:"DB_PASS"`
	Host string `env:"DB_HOST"`
	Port string `env:"DB_PORT"`
}

func Start() error {
	config := Config{}
	overenv.LoadStruct(&config)

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.Host, config.User, config.Pass, config.Name, config.Port)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	return DB.Ping()
}
