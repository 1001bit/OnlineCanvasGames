package usermodel

import (
	"context"
	"database/sql"

	"github.com/1001bit/onlinecanvasgames/services/user/pkg/crypt"
)

type User struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (us *UserStore) GetByName(ctx context.Context, username string) (*User, error) {
	user := &User{}

	err := us.db.QueryRowContext(ctx, "SELECT name, date FROM users WHERE LOWER(name) = LOWER($1)", username).Scan(&user.Name, &user.Date)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserStore) GetByNameAndPassword(ctx context.Context, username, password string) (*User, error) {
	user := &User{}
	var hash string

	// get user row regardless of character case
	row := us.db.QueryRowContext(ctx, "SELECT name, date, hash FROM users WHERE LOWER(name) = LOWER($1)", username)
	err := row.Scan(&user.Name, &user.Date, &hash)

	switch err {
	case nil:
		// no error
	case sql.ErrNoRows:
		// incorrect username
		return nil, ErrLogin
	default:
		return nil, err
	}

	// incorrect password
	if !crypt.CheckHash(password, hash) {
		return nil, ErrLogin
	}

	return user, nil
}

func (us *UserStore) Insert(ctx context.Context, username, password string) (*User, error) {
	// check existance
	var exists bool

	err := us.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(name) = LOWER($1))", username).Scan(&exists)
	if err != nil {
		return nil, err
	}

	// user already exists, can't insert a new one
	if exists {
		return nil, ErrRegister
	}

	// generate hash for user
	hash, err := crypt.GenerateHash(password)
	if err != nil {
		return nil, err
	}
	// insert into a database
	us.db.QueryRowContext(ctx, "INSERT INTO users (name, hash) VALUES ($1, $2)", username, hash)

	newUser := &User{Name: username}
	return newUser, nil
}
