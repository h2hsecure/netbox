package token

import (
	"fmt"
	"time"

	"git.h2hsecure.com/ddos/waf/internal/core/domain"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("your-secret-key")

func CreateToken(userId, ip string, validDuration time.Duration) (string, error) {
	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, domain.SessionClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(validDuration * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ddos-protector",
			Subject:   userId,
			ID:        "1",
			Audience:  []string{"somebody_else"},
		},
		UserId: userId,
		Ip:     ip,
	})

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("singing token: %w", err)
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*domain.SessionClaim, error) {
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// Check for verification errors
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	sessionClaim, ok := token.Claims.(*domain.SessionClaim)

	if !ok {
		return nil, fmt.Errorf("unknown claim type")
	}

	// Return the verified token
	return sessionClaim, nil
}
