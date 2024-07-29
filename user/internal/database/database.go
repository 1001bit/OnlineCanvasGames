package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Start() error {
	config := NewReadyConfig()

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.host, config.user, config.pass, config.name, config.port)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	return DB.Ping()
}
