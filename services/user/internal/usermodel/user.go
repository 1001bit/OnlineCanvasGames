package usermodel

import (
	"context"
	"database/sql"
	"time"

	"github.com/1001bit/onlinecanvasgames/services/user/internal/database"
	"github.com/1001bit/onlinecanvasgames/services/user/pkg/crypt"
)

const maxQueryTime = 5 * time.Second

type User struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

func NewUser() *User {
	return &User{}
}

func GetByName(ctx context.Context, username string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, maxQueryTime)
	defer cancel()

	user := NewUser()
	user.Name = username

	err := database.DB.QueryRowContext(ctx, "SELECT date FROM users WHERE name = $1", username).Scan(&user.Date)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetByNameAndPassword(ctx context.Context, username, password string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, maxQueryTime)
	defer cancel()

	user := NewUser()
	var hash string

	// get user row regardless of character case
	row := database.DB.QueryRowContext(ctx, "SELECT name, date, hash FROM users WHERE LOWER(name) = LOWER($1)", username)
	err := row.Scan(&user.Name, &user.Date, &hash)

	switch err {
	case nil:
		// no error
	case sql.ErrNoRows:
		// incorrect username
		return nil, ErrNoUser
	default:
		return nil, err
	}

	// incorrect password
	if !crypt.CheckHash(password, hash) {
		return nil, ErrNoUser
	}

	return user, nil
}

func Insert(ctx context.Context, username, password string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, maxQueryTime)
	defer cancel()

	// check existance
	var exists bool

	err := database.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(name) = LOWER($1))", username).Scan(&exists)
	if err != nil {
		return nil, err
	}

	// user already exists, can't insert a new one
	if exists {
		return nil, ErrUserExists
	}

	// create new user
	newUser := &User{Name: username}
	// generate hash for user
	hash, err := crypt.GenerateHash(password)
	if err != nil {
		return nil, err
	}
	// insert into a database
	database.DB.QueryRowContext(ctx, "INSERT INTO users (name, hash) VALUES ($1, $2)", username, hash)

	return newUser, nil
}
