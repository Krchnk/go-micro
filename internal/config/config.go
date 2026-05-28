package config

import "os"

type Config struct {
	HTTPPort string
}

func Load() Config {
	return Config{
		HTTPPort: getEnv("HTTP_PORT", "8080"),
	}
}

func (c Config) HTTPAddr() string {
	if c.HTTPPort == "" {
		return ":8080"
	}
	if c.HTTPPort[0] == ':' {
		return c.HTTPPort
	}
	return ":" + c.HTTPPort
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
