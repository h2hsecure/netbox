package service

import (
	"context"
	"fmt"

	"github.com/h2hsecure/netbox/internal/core/domain"
	"github.com/h2hsecure/netbox/internal/core/ports"
)

type serviceImpl struct {
	cache         ports.Cache
	mq            ports.MessageQueue
	token         ports.TokenService
	country       ports.CountryAdater
	cfg           domain.ConfigParams
	eventWorker   *Dispatcher
	countryPolicy map[string]domain.CountryPolicyOperation
}

func New(cache ports.Cache, mq ports.MessageQueue,
	token ports.TokenService, country ports.CountryAdater,
	cfg domain.ConfigParams) (ports.Service, error) {

	countryPolicy, err := cfg.User.CountryPolicy.Parse()
	if err != nil {
		return nil, fmt.Errorf("service building: %w", err)
	}

	return &serviceImpl{
		cache:         cache,
		mq:            mq,
		token:         token,
		country:       country,
		cfg:           cfg,
		eventWorker:   NewDispatcher(10, 100),
		countryPolicy: countryPolicy,
	}, nil
}

func (s *serviceImpl) Stop() error {
	s.eventWorker.Close()
	return nil
}

func (s *serviceImpl) putEvent(event domain.UserIpTime) {
	s.eventWorker.Push(func(ctx context.Context) error {
		return s.mq.Sent(ctx, event)
	})
}
