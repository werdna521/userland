package security

import (
	"time"

	"github.com/thanhpk/randstr"
)

const randomIDBytes = 128 / 8 // 128-bit

const TokenLife = 5 * time.Minute

type RandomID string

func GenerateRandomID() RandomID {
	return RandomID(randstr.Hex(randomIDBytes))
}
