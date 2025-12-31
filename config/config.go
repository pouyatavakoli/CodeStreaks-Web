package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseHost          string
	DatabasePort          int
	DatabaseUser          string
	DatabasePassword      string
	DatabaseName          string
	ServerPort            string
	UpdateIntervalMinutes int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	// Required variables
	required := map[string]*string{
		"DB_HOST":     &cfg.DatabaseHost,
		"DB_USER":     &cfg.DatabaseUser,
		"DB_PASSWORD": &cfg.DatabasePassword,
		"DB_NAME":     &cfg.DatabaseName,
	}

	for key, target := range required {
		value := os.Getenv(key)
		if value == "" {
			return nil, fmt.Errorf("missing required environment variable: %s", key)
		}
		*target = value
	}

	// DB_PORT (required, int)
	dbPortStr := os.Getenv("DB_PORT")
	if dbPortStr == "" {
		return nil, fmt.Errorf("missing required environment variable: DB_PORT")
	}

	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT value: %w", err)
	}
	cfg.DatabasePort = dbPort

	// Optional: SERVER_PORT (default: 8080)
	cfg.ServerPort = os.Getenv("SERVER_PORT")
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}

	// Optional: UPDATE_INTERVAL_MINUTES (default: 60)
	updateIntervalStr := os.Getenv("UPDATE_INTERVAL_MINUTES")
	if updateIntervalStr == "" {
		cfg.UpdateIntervalMinutes = 60
	} else {
		interval, err := strconv.Atoi(updateIntervalStr)
		if err != nil {
			return nil, fmt.Errorf("invalid UPDATE_INTERVAL_MINUTES value: %w", err)
		}
		cfg.UpdateIntervalMinutes = interval
	}

	return cfg, nil
}
