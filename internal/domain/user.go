package domain

import "time"

type User struct {
	ID                 int64     `json:"id" db:"id"`
	Handle             string    `json:"handle" db:"handle"`
	CurrentStreak      int       `json:"current_streak" db:"current_streak"`
	MaxStreak          int       `json:"max_streak" db:"max_streak"`
	LastSubmissionDate time.Time `json:"last_submission_date" db:"last_submission_date"`
	LastUpdatedAt      time.Time `json:"last_updated_at" db:"last_updated_at"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}
