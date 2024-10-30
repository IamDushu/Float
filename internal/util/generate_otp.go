package util

import (
	"crypto/rand"
	"math/big"
)

func GenerateOTP() (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(90000))
	if err != nil {
		return 0, err
	}

	return int(n.Int64()) + 10000, nil
}
