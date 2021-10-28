package jwt

import "time"

const (
	AccessTokenLife = 5 * time.Minute
	// TODO: make expiration longer, current duration is just for debugging purposes
	RefreshTokenLife = 10 * time.Minute
)
