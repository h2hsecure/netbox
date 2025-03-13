package token

import (
	"fmt"
	"time"

	"git.h2hsecure.com/ddos/waf/internal/core/domain"
	"git.h2hsecure.com/ddos/waf/internal/core/ports"
	"github.com/golang-jwt/jwt/v5"
)

type tokenBuilder struct {
	secretKey            []byte
	defaultValidDuration time.Duration
}

func NewTokenService(secretKey string, defaultValidDuration time.Duration) ports.TokenService {
	return &tokenBuilder{
		defaultValidDuration: defaultValidDuration,
		secretKey:            []byte(secretKey),
	}
}

func (t *tokenBuilder) CreateToken(userId, ip string, validDuration time.Duration) (string, error) {

	if validDuration == 0 {
		validDuration = t.defaultValidDuration
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &domain.SessionClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(validDuration)),
			Subject:   userId,
		},
	})

	tokenString, err := claims.SignedString(t.secretKey)
	if err != nil {
		return "", fmt.Errorf("singing token: %w", err)
	}

	return tokenString, nil
}

func (t *tokenBuilder) VerifyToken(tokenString string) (*domain.SessionClaim, error) {
	// Parse the token with the secret key
	token, err := jwt.ParseWithClaims(tokenString, &domain.SessionClaim{}, func(token *jwt.Token) (any, error) {
		return t.secretKey, nil
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
