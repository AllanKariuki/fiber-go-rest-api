package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl string
	Port string
	JWTSecret string
}

var AppConfig *Config

func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		DBUrl: getEnv("DATABASE_URL", ""),
		Port: getEnv("PORT", "8081"),
		JWTSecret: getEnv("JWT_SECRET", ""),
	}
}

func getEnv (key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}