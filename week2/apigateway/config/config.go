package config

import (
	"os"
	"go.uber.org/zap"
)

type Config struct {
	GatewayPort string
	Service1URL string
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

	return cfg, nil
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
