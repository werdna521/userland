package security

import "github.com/thanhpk/randstr"

const randomIDBytes = 128 / 8 // 128-bit

type RandomID string

func GenerateRandomID() RandomID {
	return RandomID(randstr.Hex(randomIDBytes))
}
