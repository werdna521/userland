package repository

import "time"

type PasswordHistory struct {
	ID        string
	UserID    string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
