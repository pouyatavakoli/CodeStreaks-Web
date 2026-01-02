package service

import (
	"errors"
	"log"
	"time"

	"github.com/pouyatavakoli/CodeStreaks-web/internal/domain"
	"github.com/pouyatavakoli/CodeStreaks-web/internal/repository"
	"github.com/pouyatavakoli/CodeStreaks-web/pkg/codeforces"
	"gorm.io/gorm"
)

type UserService interface {
	AddUser(handle string) (*domain.User, error)
	GetLeaderboard(page, pageSize int) ([]domain.UserResponse, int64, error)
	GetUserByHandle(handle string) (*domain.User, error)
	UpdateUserStreaks(user *domain.User, submissions []domain.CodeforcesSubmission) error
}

type userService struct {
	userRepo       repository.UserRepository
	submissionRepo repository.SubmissionRepository
	cfClient       *codeforces.Client
}

func NewUserService(
	userRepo repository.UserRepository,
	submissionRepo repository.SubmissionRepository,
	cfClient *codeforces.Client,
) UserService {
	return &userService{
		userRepo:       userRepo,
		submissionRepo: submissionRepo,
		cfClient:       cfClient,
	}
}

func (s *userService) AddUser(handle string) (*domain.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByHandle(handle)
	if err == nil {
		return existingUser, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Validate handle with Codeforces API
	userInfo, err := s.cfClient.GetUserInfo(handle)
	if err != nil {
		return nil, err
	}

	// Create new user
	user := &domain.User{
		CodeforcesHandle: userInfo.Handle,
		Rating:           userInfo.Rating,
		Rank:             userInfo.Rank,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	log.Printf("Added new user: %s", handle)
	return user, nil
}

func (s *userService) GetLeaderboard(page, pageSize int) ([]domain.UserResponse, int64, error) {
	offset := (page - 1) * pageSize

	users, err := s.userRepo.GetLeaderboard(pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.userRepo.CountUsers()
	if err != nil {
		return nil, 0, err
	}

	responses := make([]domain.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse(offset + i + 1)
	}

	return responses, total, nil
}

func (s *userService) GetUserByHandle(handle string) (*domain.User, error) {
	return s.userRepo.FindByHandle(handle)
}

func (s *userService) UpdateUserStreaks(user *domain.User, submissions []domain.CodeforcesSubmission) error {
	if len(submissions) == 0 {
		return nil
	}

	// Calculate streak
	streak := s.calculateStreak(submissions)

	// Update user fields
	user.CurrentStreak = streak
	if streak > user.MaxStreak {
		user.MaxStreak = streak
	}

	latestSubmission := submissions[0]
	submissionTime := time.Unix(latestSubmission.CreationTimeSeconds, 0)
	user.LastSubmissionAt = &submissionTime

	return s.userRepo.Update(user)
}

func (s *userService) calculateStreak(submissions []domain.CodeforcesSubmission) int {
	if len(submissions) == 0 {
		return 0
	}

	// Group submissions by date
	submissionsByDate := make(map[string]bool)
	for _, sub := range submissions {
		if sub.Verdict == "OK" {
			date := time.Unix(sub.CreationTimeSeconds, 0).Format("2006-01-02")
			submissionsByDate[date] = true
		}
	}

	// Calculate consecutive days
	streak := 0
	currentDate := time.Now()

	for {
		dateStr := currentDate.Format("2006-01-02")
		if submissionsByDate[dateStr] {
			streak++
			currentDate = currentDate.AddDate(0, 0, -1)
		} else {
			// Allow one day gap for today
			if streak == 0 && dateStr == time.Now().Format("2006-01-02") {
				currentDate = currentDate.AddDate(0, 0, -1)
				continue
			}
			break
		}
	}

	return streak
}
