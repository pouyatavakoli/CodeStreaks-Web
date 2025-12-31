package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/pouyatavakoli/CodeStreaks-web/config"
)

func SetupDatabase() (*sql.DB, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	db, err := NewPostgresDBFromAppConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	if err := TestConnection(db); err != nil {
		return nil, fmt.Errorf("database connection test failed: %w", err)
	}

	log.Println("Database connection established successfully")
	return db, nil
}

func MonitorConnectionPool(db *sql.DB, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		stats := Stats(db)
		log.Printf("Database Pool Stats - OpenConnections: %d, InUse: %d, Idle: %d",
			stats.OpenConnections,
			stats.InUse,
			stats.Idle,
		)
	}
}
