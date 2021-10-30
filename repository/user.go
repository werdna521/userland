package repository

import (
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
