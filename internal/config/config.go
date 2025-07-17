package config

import (
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Scrapper ScrapperConfig
}

type DatabaseConfig struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
}

type ServerConfig struct {
	Port string
}

type ScrapperConfig struct {
	BaseURL string
	Timeout int
}

func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			DSN:          getEnv("DATABASE_DSN", ""),
			MaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 5),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Scrapper: ScrapperConfig{
			BaseURL: getEnv("SCRAPPER_BASE_URL", "https://resultadodelaloteria.com"),
			Timeout: getEnvInt("SCRAPPER_TIMEOUT", 30),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
