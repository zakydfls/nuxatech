package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	JwtAccessSecret  string
	JwtRefreshSecret string
	Port             string
	AccessTokenTTL   int
	RefreshTokenTTL  int
	DbName           string
	DbHost           string
	DbPort           string
	DbUser           string
	DbPass           string
}

var Envs = InitConfig()

func InitConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	return &Config{
		JwtAccessSecret:  getEnv("JWT_ACCESS_SECRET", "secret"),
		JwtRefreshSecret: getEnv("JWT_REFRESH_SECRET", "refresh"),
		Port:             getEnv("PORT", "3000"),
		AccessTokenTTL:   getEnvAsInt("ACC_EXPIRED", 3600*3),
		RefreshTokenTTL:  getEnvAsInt("REFRESH_EXPIRED", 3600*24*7),
		DbName:           getEnv("DB_NAME", "postgres"),
		DbHost:           getEnv("DB_HOST", "127.0.0.1"),
		DbPort:           getEnv("DB_PORT", "5432"),
		DbUser:           getEnv("DB_USER", "postgres"),
		DbPass:           getEnv("DB_PASS", "postgres"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvAsInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return intValue
}
