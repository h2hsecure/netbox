package service

import (
	"context"
	"errors"

	"git.h2hsecure.com/ddos/waf/internal/core/domain"
	"github.com/rs/zerolog/log"
)

func (s *serviceImpl) AccessAtempt(ctx context.Context, token string, event domain.AttemptRequest) (op domain.AttemptOperation) {
	op = domain.AttemptUserAllow

	// check processing from begging
	if s.cfg.User.DisableProcessing {
		return
	}

	defer func() {
		currentCounter, err := s.cache.Inc(ctx, "c"+event.Ip, 1)

		if err != nil {
			log.Err(err).
				Interface("event", event).
				Msg("inc ip counter")
		}

		if currentCounter%s.cfg.User.CounterFreq == 0 {
			s.putEvent(event.UserIpTime)
		}
	}()

	// check location from request
	if event.Location != nil {
		cop, has := s.countryPolicy[*event.Location]

		if !has {
			cop = s.countryPolicy[domain.CountryPolicyAll]
		}

		switch cop {
		case domain.CountryPolicyOperationNoop:
			return
		case domain.CountryPolicyOperationDeny:
			op = domain.AttemptDenyUserByCountry
			return
		}
	}

	// ip based control in cache
	last, err := s.cache.Get(ctx, event.Ip)

	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		log.Err(err).
			Interface("event", event).
			Msg("fetching ip from cache")
	}

	if last != "" {
		log.Warn().
			Str("ip", last).
			Msg("ip found in cache")
		op = domain.AttemptDenyUserByIp

		return
	}

	// verify token
	t, err := s.token.VerifyToken(token)

	if err != nil {
		log.Err(err).Send()
		op = domain.AttemptValidate
		return
	}

	// verify user in cache
	last, err = s.cache.Get(ctx, t.Subject)

	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		log.Err(err).
			Interface("event", event).
			Msg("feching user from cache")
	}

	if last != "" {
		log.Warn().Str("sub", last).Msg("user found in cache")
		op = domain.AttemptValidate
		return
	}

	event.User = t.Subject

	currentCounter, err := s.cache.Inc(ctx, "c"+event.User, 1)

	if err != nil {
		log.Err(err).
			Interface("event", event).
			Msg("inc user counter")
	}

	if currentCounter%s.cfg.User.CounterFreq == 0 {
		s.putEvent(event.UserIpTime)
	}

	return
}
