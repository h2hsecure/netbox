package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/h2hsecure/netbox/internal/core/domain"
	"github.com/h2hsecure/netbox/internal/core/ports"
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

	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) && !errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		return nil, fmt.Errorf("verfiy token: %w", err)
	}

	// Check if the token is expried
	if errors.Is(err, jwt.ErrTokenExpired) {
		return nil, domain.ErrTokenExperied
	}

	if !token.Valid {
		return nil, domain.ErrTokenInvalid
	}

	sessionClaim, ok := token.Claims.(*domain.SessionClaim)

	if !ok {
		return nil, fmt.Errorf("wrong claim: %w", domain.ErrTokenInvalid)
	}

	// Return the verified token
	return sessionClaim, nil
}
