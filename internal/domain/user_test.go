package domain

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUser_JSONMarshal(t *testing.T) {
	user := User{
		ID:                 123,
		Handle:             "john_doe",
		CurrentStreak:      5,
		MaxStreak:          10,
		LastSubmissionDate: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		LastUpdatedAt:      time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC),
		CreatedAt:          time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal User to JSON: %v", err)
	}

	t.Logf("JSON output: %s", string(jsonBytes))

	// Unmarshal back to verify
	var decodedUser User
	err = json.Unmarshal(jsonBytes, &decodedUser)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON back to User: %v", err)
	}

	// Verify all fields match
	if decodedUser.ID != user.ID {
		t.Errorf("ID mismatch: got %d, want %d", decodedUser.ID, user.ID)
	}

	if decodedUser.Handle != user.Handle {
		t.Errorf("Handle mismatch: got %s, want %s", decodedUser.Handle, user.Handle)
	}

	if decodedUser.CurrentStreak != user.CurrentStreak {
		t.Errorf("CurrentStreak mismatch: got %d, want %d", decodedUser.CurrentStreak, user.CurrentStreak)
	}

	if decodedUser.MaxStreak != user.MaxStreak {
		t.Errorf("MaxStreak mismatch: got %d, want %d", decodedUser.MaxStreak, user.MaxStreak)
	}

	if !decodedUser.LastSubmissionDate.Equal(user.LastSubmissionDate) {
		t.Errorf("LastSubmissionDate mismatch")
	}

	if !decodedUser.LastUpdatedAt.Equal(user.LastUpdatedAt) {
		t.Errorf("LastUpdatedAt mismatch")
	}

	if !decodedUser.CreatedAt.Equal(user.CreatedAt) {
		t.Errorf("CreatedAt mismatch")
	}

	t.Log("Test passed: User successfully marshaled and unmarshaled from JSON")
}
