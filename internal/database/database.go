package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Start() error {
	config := NewConfig()

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.user, config.pass, config.name)

	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	return DB.Ping()
}

// 2006-01-02T15:04:05Z -> 2 January 2006
func FormatPostgresDate(dateStr string) (string, error) {
	t, err := time.Parse("2006-01-02T15:04:05Z", dateStr)
	if err != nil {
		return "", err
	}

	return t.Format("2 January 2006"), nil
}
