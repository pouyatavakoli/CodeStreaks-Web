package domain

import (
	"time"
)

type User struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	CodeforcesHandle string     `gorm:"uniqueIndex;not null" json:"codeforces_handle"`
	CurrentStreak    int        `gorm:"default:0" json:"current_streak"`
	MaxStreak        int        `gorm:"default:0" json:"max_streak"`
	LastSubmissionAt *time.Time `json:"last_submission_at"`
	Rating           int        `gorm:"default:0" json:"rating"`
	Rank             string     `json:"rank"`
	TotalSubmissions int        `gorm:"default:0" json:"total_submissions"`
	IsActive         bool       `gorm:"default:true" json:"is_active"`
	LastCheckedAt    *time.Time `json:"last_checked_at"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type UserResponse struct {
	ID               uint       `json:"id"`
	CodeforcesHandle string     `json:"codeforces_handle"`
	CurrentStreak    int        `json:"current_streak"`
	MaxStreak        int        `json:"max_streak"`
	LastSubmissionAt *time.Time `json:"last_submission_at"`
	Rating           int        `json:"rating"`
	Rank             string     `json:"rank"`
	TotalSubmissions int        `json:"total_submissions"`
	LeaderboardRank  int        `json:"leaderboard_rank"`
}

func (u *User) ToResponse(rank int) UserResponse {
	return UserResponse{
		ID:               u.ID,
		CodeforcesHandle: u.CodeforcesHandle,
		CurrentStreak:    u.CurrentStreak,
		MaxStreak:        u.MaxStreak,
		LastSubmissionAt: u.LastSubmissionAt,
		Rating:           u.Rating,
		Rank:             u.Rank,
		LeaderboardRank:  rank,
	}
}
