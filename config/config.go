package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database   DatabaseConfig
	Server     ServerConfig
	Codeforces CodeforcesConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ServerConfig struct {
	Port string
	Env  string
}

type CodeforcesConfig struct {
	BaseURL        string
	WorkerPoolSize int
	UpdateInterval int // seconds
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	workerPoolSize, _ := strconv.Atoi(getEnv("WORKER_POOL_SIZE", "10"))
	updateInterval, _ := strconv.Atoi(getEnv("UPDATE_INTERVAL", "60"))

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "CodeStreaks"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Codeforces: CodeforcesConfig{
			BaseURL:        getEnv("CODEFORCES_API_URL", "https://codeforces.com/api"),
			WorkerPoolSize: workerPoolSize,
			UpdateInterval: updateInterval,
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
