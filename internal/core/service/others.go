package service

import (
	"context"
	"net/netip"

	"git.h2hsecure.com/ddos/waf/internal/core/domain"
)

func (s *serviceImpl) ManageIps(context.Context, []netip.Addr, domain.IpOperation) error {
	return nil
}
func (s *serviceImpl) QueryIp(context.Context, netip.Addr) (domain.IpOperation, error) {
	return domain.IpOperationAllow, nil
}

func (s *serviceImpl) CurrentConfig(context.Context) domain.ConfigParams {
	return domain.ConfigParams{}
}
func (s *serviceImpl) Health(context.Context) error {
	return nil
}
