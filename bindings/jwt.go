package bindings

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthClaims struct {
	OriginalLoginAt *jwt.NumericDate
	jwt.RegisteredClaims
}

func generateAccessToken(subject string, expirationSeconds int, issuer string, originalLoginTime time.Time, secretKey string) (string, error) {
	claims := AuthClaims{
		jwt.NewNumericDate(originalLoginTime),
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
			Subject:   subject,
		},
	}

	if expirationSeconds > 0 {
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(expirationSeconds) * time.Second))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	return signedToken, err
}

func verifyAccessToken(signedToken string, secretKey string, expectedIssuer string) (AuthClaims, error) {
	token, err := jwt.ParseWithClaims(signedToken, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	}, jwt.WithIssuer(expectedIssuer))

	if err != nil {
		return AuthClaims{}, err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok {
		return *claims, nil
	}

	return AuthClaims{}, errors.New("bad claims")
}
