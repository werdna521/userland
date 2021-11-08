package repository

import "time"

// TODO: store IP as well
type Session struct {
	ID        string
	UserID    string
	Client    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AccessToken struct {
	ID        string
	SessionID string
	UserID    string
}

type RefreshToken struct {
	ID        string
	SessionID string
	UserID    string
}
