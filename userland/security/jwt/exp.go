package jwt

import "time"

const (
	AccessTokenLife  = 1 * time.Hour
	RefreshTokenLife = 24 * time.Hour
)
