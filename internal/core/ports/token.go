package ports

import (
	"time"

	"git.h2hsecure.com/ddos/waf/internal/core/domain"
)

type TokenService interface {
	CreateToken(userId, ip string, validDuration time.Duration) (string, error)
	VerifyToken(tokenString string) (*domain.SessionClaim, error)
}
