package usermodel

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/1001bit/OnlineCanvasGames/internal/crypt"
	"github.com/1001bit/OnlineCanvasGames/internal/database"
)

const maxQueryTime = 5 * time.Second

var (
	ErrNoUserExists = errors.New("user with such name doesn't exist")
	ErrNoSuchUser   = errors.New("incorrect username or password")
	ErrUserExists   = errors.New("user with such name already exists")
)

type User struct {
	ID   int
	Name string
	Date string
}

func NewUser() *User {
	return &User{}
}

func GetByID(ctx context.Context, userID int) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, maxQueryTime)
	defer cancel()

	user := NewUser()
	user.ID = userID

	err := database.DB.QueryRowContext(ctx, "SELECT name, date FROM users WHERE id = $1", userID).Scan(&user.Name, &user.Date)
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

	row := database.DB.QueryRowContext(ctx, "SELECT id, name, date, hash FROM users WHERE LOWER(name) = LOWER($1)", username)
	err := row.Scan(&user.ID, &user.Name, &user.Date, &hash)

	switch err {
	case nil:
		// no error
	case sql.ErrNoRows:
		return nil, ErrNoSuchUser
	default:
		return nil, err
	}

	if !crypt.CheckHash(password, hash) {
		return nil, ErrNoSuchUser
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

	if exists {
		return nil, ErrUserExists
	}

	// create new user
	newUser := &User{Name: username}

	hash, err := crypt.GenerateHash(password)
	if err != nil {
		return nil, err
	}

	err = database.DB.QueryRowContext(ctx, "INSERT INTO users (name, hash) VALUES ($1, $2) RETURNING id", username, hash).Scan(&newUser.ID)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func IDExists(ctx context.Context, userID int) bool {
	ctx, cancel := context.WithTimeout(ctx, maxQueryTime)
	defer cancel()

	// check existance
	exists := false

	err := database.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userID).Scan(&exists)
	if err != nil {
		return false
	}

	return exists
}
