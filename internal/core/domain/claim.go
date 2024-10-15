package domain

import "github.com/golang-jwt/jwt/v5"

type SessionClaim struct {
	jwt.RegisteredClaims
	UserId string `json:"userId"`
	Ip     string `json:"ip"`
}

func WithDefaultCliam(userId, ip string, claim jwt.RegisteredClaims) SessionClaim {
	return SessionClaim{
		claim,
		userId,
		ip,
	}
}
