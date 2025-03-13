package cmd

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type CacheParams struct {
	Sock    string `env:"CACHE_SOCK"`
	Size    int    `env:"CACHE_SIZE"`
	MaxUser uint64 `env:"MAX_USER"`
	MaxIp   uint64 `env:"MAX_IP"`
	MaxPath uint64 `env:"MAX_PATH"`
}

type EnforcerParams struct {
	ClusterStr string `env:"CLUSTER_STR" description:"Raft Cluster string for cluster itself. comman delimeted string. cluster format is host_id:host_name:raft_port:grpc_port"`
	MyAddress  string `env:"MY_ADDRESS" description:"my address in raft cluster"`
}
type NginxParams struct {
	ContextPath  string `env:"CONTEXT_PATH"`
	BackendHost  string `env:"BACKEND_HOST"`
	BackendPort  string `env:"BACKEND_PORT"`
	Domain       string `env:"DOMAIN"`
	DomainProto  string `env:"DOMAIN_PROTO"`
	InternalSock string `env:"INTERNAL_SOCK"`
}

type UserParams struct {
	DefaultLocale     string        `env:"DEFAULT_LOCALE"`
	BackendLogo       string        `env:"BACKEND_LOGO"`
	TokenDuration     time.Duration `env:"TOKEN_DURATION"`
	CookieDuration    time.Duration `env:"COOKIE_DURATION"`
	CookieName        string        `env:"COOKIE_NAME"`
	TokenSecret       string        `env:"TOKEN_SECRET"`
	InsecureCookie    bool          `env:"ENABLE_INSECURE_COOCKIE"`
	ChallengeHmac     string        `env:"CHALLENGE_HMAC_KEY"`
	SearchEngineBots  bool          `env:"ENABLE_SEARCH_ENGINE_BOTS"`
	DisableProcessing bool          `env:"DISABLE_PROCESSING"`
}

type ConfigParams struct {
	Nginx    NginxParams
	Enforcer EnforcerParams
	Cache    CacheParams
	User     UserParams

	LogDir       string `env:"LOG_DIR"`
	PromListener string `env:"PROM_LISTEN"`
	SystemId     string `env:"SYSTEM_ID"`
}

func CurrentConfig() (ConfigParams, error) {

	var cfg ConfigParams

	if err := env.Parse(&cfg); err != nil {
		return ConfigParams{}, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, nil
}
