package repository

import (
	"database/sql"
	"time"
)

type User struct {
	ID        string
	Email     string
	Password  string
	IsActive  bool
	UserBio   *UserBio
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

type UserBio struct {
	ID        string
	Fullname  string
	Location  string
	Bio       string
	Web       string
	Picture   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
