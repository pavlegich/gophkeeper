// Package hash contains objects and methods for creating
// and validating tokens.
package hash

import (
	"context"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims contains objects for claiming JWT
// and catching ID.
type Claims struct {
	jwt.RegisteredClaims
	ID int
}

// Token contains objects for signing tokens.
type Token struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	tokenExp   time.Duration
}

// NewToken returns new token object.
func NewToken(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, tokenExp time.Duration) *Token {
	return &Token{
		privateKey: privateKey,
		publicKey:  publicKey,
		tokenExp:   tokenExp,
	}
}

// Create creates token and returns it as a string.
func (t *Token) Create(ctx context.Context, id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.tokenExp)),
		},
		ID: id,
	})

	tokenString, err := token.SignedString(t.privateKey)
	if err != nil {
		return "", fmt.Errorf("Create: sign string with key failed %w", err)
	}

	return tokenString, nil
}

// Validate returns received data from the token for authentication.
func (t *Token) Validate(tokenString string) (int, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(j *jwt.Token) (interface{}, error) {
			if _, ok := j.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Validate: unexpected signing method: %v", j.Header["alg"])
			}
			return t.publicKey, nil
		})
	if err != nil {
		return -1, fmt.Errorf("Validate: parse token failed %w", err)
	}
	if !token.Valid {
		return -1, fmt.Errorf("Validate: token is not valid")
	}

	return claims.ID, nil
}
