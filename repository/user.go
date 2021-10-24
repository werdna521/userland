package repository

import (
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
