package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Sub string
	Iat time.Time
	Nbf time.Time
	Exp time.Time
}

func CreateUnsignedJWT(c Claims) (string, error) {
	claims := jwt.MapClaims{
		"sub": c.Sub,
		"iat": c.Iat,
		"nbf": c.Nbf,
		"exp": c.Exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)

	unsignedToken, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		return "", fmt.Errorf("error creating token: %w", err)
	}

	return unsignedToken, nil
}
