package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Database    *DatabaseConfig
	Environment string
	ServerPort  int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  bool
}

func GetConfig() *Config {
	cfg := &Config{
		Database: &DatabaseConfig{
			Host:     getEnvOrDefault("DATABASE_HOST", "localhost"),
			User:     getEnvOrDefault("DATABASE_USER", "admin"),
			Password: getEnvOrDefault("DATABASE_PASSWORD", "admin"),
			Name:     getEnvOrDefault("DATABASE_NAME", "mike-local-db"),
			Port:     getEnvOrDefault("DATABASE_PORT", 5411),
			SSLMode:  getEnvOrDefault("DATABASE_SSL_MODE", false),
		},
		ServerPort: getEnvOrDefault("SERVER_PORT", 8080),
	}
	switch strings.ToLower(os.Getenv("ENV")) {
	case "development":
		loadDevelopmentConfig(cfg)
	default:
		loadDevelopmentConfig(cfg)
	}
	return cfg
}

func getEnvOrDefault[T string | int | bool](envVarName string, defaultVal T) T {
	val := os.Getenv(envVarName)
	if val == "" {
		return defaultVal
	}
	switch any(defaultVal).(type) {
	case string:
		return any(val).(T)
	case int:
		i, _ := strconv.Atoi(val)
		// don't error check cause we WANT it to blow up if it's not a parseable int
		return any(i).(T)
	case bool:
		b, _ := strconv.ParseBool(val)
		// don't error check cause we WANT it to blow up if it's not a parseable bool
		return any(b).(T)
	}
	return defaultVal
}
