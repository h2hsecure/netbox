package service_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"git.h2hsecure.com/ddos/waf/cmd"
	"git.h2hsecure.com/ddos/waf/internal/core/ports"
	"git.h2hsecure.com/ddos/waf/internal/core/service"
	"git.h2hsecure.com/ddos/waf/internal/repository/cache"
	"git.h2hsecure.com/ddos/waf/internal/repository/grpc"
	"git.h2hsecure.com/ddos/waf/internal/repository/token"
)

const CountryPolicyStr = "TR:noop,EN:allow,NL:deny,*:allow"

var (
	testService ports.Service

	testCtx      = context.Background()
	mockCache    *cache.MockCache
	tokenService ports.TokenService
)

func TestMain(m *testing.M) {

	cfg, _ := cmd.CurrentConfig()

	cfg.User.CountryPolicy = CountryPolicyStr
	cfg.User.CounterFreq = 1
	cfg.User.TokenSecret = "02b2648e69c33bac085b"

	mockCache = cache.CreateMockCache()
	tokenService = token.NewTokenService(cfg.User.TokenSecret, 1*time.Hour)

	var err error
	testService, err = service.New(mockCache, grpc.NewMockMq(), tokenService, nil, cfg)

	if err != nil {
		fmt.Printf("service init error: %v", err)
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}
