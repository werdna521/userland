package repository

import (
	"context"
	"time"
)

type User struct {
	ID        int64
	Fullname  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
}
