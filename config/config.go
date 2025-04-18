package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
}

func LoadConfig() *config {
	err := godotenv.Load()
	if err != nil {
		log.Println("error loading .env file")
	}
	return &config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "taskdb"),
		JWTSecret:  getEnv("JWT_SECRET", "secret"),
	}

}

func getEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}
