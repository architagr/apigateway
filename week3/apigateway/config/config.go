package config

import (
	"os"
	"strconv"

	"go.uber.org/zap"
)

type Config struct {
	RateLimitter rateLimiterConfig
	GatewayPort  string
	Service1URL  string
}
type RateLimiterType string

const (
	RateLimiterType_TokenBucket       RateLimiterType = "tokenBucket"
	RateLimiterType_LeakyBucket       RateLimiterType = "leakyBucket"
	RateLimiterType_FixedWindowBucket RateLimiterType = "fixedWindow"
)

type rateLimiterConfig struct {
	LimiterType RateLimiterType
	Rate        int
	Capacity    int
}

func NewConfig() (*Config, error) {
	//default values
	cfg := &Config{
		GatewayPort: "8080",
		Service1URL: "localhost:50051",
	}

	if port := os.Getenv("GATEWAY_PORT"); port != "" {
		cfg.GatewayPort = port
	}
	if addr := os.Getenv("SERVICE1_URL"); addr != "" {
		cfg.Service1URL = addr
	}
	cfg.RateLimitter = getRateLimiterConfig()
	return cfg, nil
}

func getRateLimiterConfig() rateLimiterConfig {
	c := rateLimiterConfig{
		LimiterType: RateLimiterType_TokenBucket,
		Rate:        4,
		Capacity:    4,
	}
	if rlType := os.Getenv("RATE_LIMITTER_TYPE"); len(rlType) > 0 {
		c.LimiterType = RateLimiterType(rlType)
	}
	if rate := os.Getenv("RATE_LIMITTER_RATE"); len(rate) > 0 {
		c.Rate, _ = strconv.Atoi(rate)
	}
	if cap := os.Getenv("RATE_LIMITTER_CAPACITY"); len(cap) > 0 {
		c.Capacity, _ = strconv.Atoi(cap)
	}
	return c
}

func (c *Config) ZapConfig() zap.Config {
	return zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}
