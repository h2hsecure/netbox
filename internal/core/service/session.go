package service

import (
	"context"
	"fmt"
	"time"

	"git.h2hsecure.com/ddos/waf/internal/core/domain"
	"github.com/google/uuid"
)

func (s *serviceImpl) OpenSession(ctx context.Context, userIpTime domain.UserIpTime) (string, error) {

	userId := uuid.NewString()

	token, err := s.token.CreateToken(userId, userIpTime.Ip, time.Duration(0))

	userIpTime.User = userId
	userIpTime.Timestamp = time.Now().Unix()

	defer s.putEvent(userIpTime)

	if err != nil {
		return "", fmt.Errorf("openSession: %w", err)
	}

	return token, err
}
