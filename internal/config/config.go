package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPCPort     string
	GRPCClient   string
	MetricsAddr  string
	KafkaBroker  string
	KafkaTopic   string
	KafkaGroup   string
	JWTSecret    string
	JWTTTLHours  int
	AuthPassword string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func Load() Config {
	// Load .env if present. Environment variables from OS still have higher priority.
	_ = godotenv.Load()

	return Config{
		GRPCPort:     getEnv("GRPC_PORT", "9090"),
		GRPCClient:   getEnv("GRPC_ADDR", "localhost:9090"),
		MetricsAddr:  getEnv("METRICS_ADDR", ":2112"),
		KafkaBroker:  getEnv("KAFKA_BROKER", "localhost:9092"),
		KafkaTopic:   getEnv("KAFKA_TOPIC", "user-registered"),
		KafkaGroup:   getEnv("KAFKA_GROUP", "notify-service"),
		JWTSecret:    getEnv("JWT_SECRET", "supersecretkey"),
		JWTTTLHours:  getEnvInt("JWT_TTL_HOURS", 24),
		AuthPassword: getEnv("AUTH_PASSWORD", "password"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "go_micro"),
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

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return parsed
}
