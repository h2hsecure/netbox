package ports

import (
	"time"

	"github.com/h2hsecure/netbox/internal/core/domain"
)

type TokenService interface {
	CreateToken(userId, ip string, validDuration time.Duration) (string, error)
	VerifyToken(tokenString string) (*domain.SessionClaim, error)
}
