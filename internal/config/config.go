package config

import "os"

type Config struct {
	GRPCPort string
}

func Load() Config {
	return Config{
		GRPCPort: getEnv("GRPC_PORT", "9090"),
	}
}

func (c Config) GRPCAddr() string {
	if c.GRPCPort == "" {
		return ":9090"
	}
	if c.GRPCPort[0] == ':' {
		return c.GRPCPort
	}
	return ":" + c.GRPCPort
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
