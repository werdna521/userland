package repository

import "time"

type ForgotPassword struct {
	ID          int64
	UserID      int64
	OldPassword string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
