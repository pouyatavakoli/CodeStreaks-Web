package database

import (
	"context"
	"testing"
	"time"

	"github.com/pouyatavakoli/CodeStreaks-web/config"
)

func TestNewPostgresDB(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database integration test")
	}

	testConfig := &config.Config{
		DatabaseHost:          "localhost",
		DatabasePort:          5432,
		DatabaseUser:          "postgres",
		DatabasePassword:      "postgres",
		DatabaseName:          "test_db",
		ServerPort:            "8080",
		UpdateIntervalMinutes: 60,
	}

	db, err := NewPostgresDBFromAppConfig(testConfig)
	if err != nil {
		t.Skipf("Skipping test: could not connect to database: %v", err)
	}
	defer CloseDB(db)

	if err := TestConnection(db); err != nil {
		t.Errorf("TestConnection failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := HealthCheckWithContext(ctx, db); err != nil {
		t.Errorf("HealthCheck failed: %v", err)
	}

	stats := Stats(db)
	if stats.MaxOpenConnections != 25 {
		t.Errorf("Expected MaxOpenConnections=25, got %d", stats.MaxOpenConnections)
	}
}
