package domain

import "time"

type User struct {
	ID                 int64     `json:"id" db:"id"`
	Handle             string    `json:"handle" db:"handle"`
	CurrentStreak      int       `json:"currentStreak" db:"currentStreak"`
	MaxStreak          int       `json:"maxStreak" db:"maxStreak"`
	LastSubmissionDate time.Time `json:"lastSubmissionDate" db:"lastSubmissionDate"`
	LastUpdatedAt      time.Time `json:"lastUpdatedAt" db:"lastUpdatedAt"`
	CreatedAt          time.Time `json:"createdAt" db:"createdAt"`
}
