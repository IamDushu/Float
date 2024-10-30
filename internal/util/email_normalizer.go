package util

import (
	normalizer "github.com/dimuska139/go-email-normalizer/v3"
)

func NormalizeEmail(email string) string {
	n := normalizer.NewNormalizer()
	return n.Normalize(email)
}
