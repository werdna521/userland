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
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepository interface {
	PrepareStatements(context.Context) error
	TearDownStatements()
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUserActivationStatusByEmail(
		ctx context.Context,
		email string,
		isActive bool,
	) error
}
