package service

import (
	"log"
	"sync"
	"time"

	"github.com/pouyatavakoli/CodeStreaks-web/internal/domain"
	"github.com/pouyatavakoli/CodeStreaks-web/internal/repository"
	"github.com/pouyatavakoli/CodeStreaks-web/pkg/codeforces"
)

type SyncService interface {
	SyncAllUsers() error
	SyncUser(user *domain.User) error
}

type syncService struct {
	userRepo       repository.UserRepository
	submissionRepo repository.SubmissionRepository
	cfClient       *codeforces.Client
	workerPoolSize int
}

func NewSyncService(
	userRepo repository.UserRepository,
	submissionRepo repository.SubmissionRepository,
	cfClient *codeforces.Client,
	workerPoolSize int,
) SyncService {
	return &syncService{
		userRepo:       userRepo,
		submissionRepo: submissionRepo,
		cfClient:       cfClient,
		workerPoolSize: workerPoolSize,
	}
}

type syncJob struct {
	user *domain.User
}

type syncResult struct {
	user *domain.User
	err  error
}

func (s *syncService) SyncAllUsers() error {
	users, err := s.userRepo.GetAllActiveUsers()
	if err != nil {
		return err
	}

	if len(users) == 0 {
		log.Println("No active users to sync")
		return nil
	}

	log.Printf("Starting sync for %d users with %d workers", len(users), s.workerPoolSize)
	startTime := time.Now()

	// Create channels
	jobs := make(chan syncJob, len(users))
	results := make(chan syncResult, len(users))

	// Start worker pool
	var wg sync.WaitGroup
	for w := 0; w < s.workerPoolSize; w++ {
		wg.Add(1)
		go s.worker(w, jobs, results, &wg)
	}

	// Send jobs
	for i := range users {
		jobs <- syncJob{user: &users[i]}
	}
	close(jobs)

	// Wait for workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	successCount := 0
	errorCount := 0
	for result := range results {
		if result.err != nil {
			log.Printf("Error syncing user %s: %v", result.user.CodeforcesHandle, result.err)
			errorCount++
		} else {
			successCount++
		}
	}

	duration := time.Since(startTime)
	log.Printf("Sync completed: %d successful, %d errors in %v", successCount, errorCount, duration)

	return nil
}

func (s *syncService) worker(id int, jobs <-chan syncJob, results chan<- syncResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		// Add small delay to respect API rate limits
		time.Sleep(200 * time.Millisecond)

		err := s.SyncUser(job.user)
		results <- syncResult{
			user: job.user,
			err:  err,
		}
	}
}

func (s *syncService) SyncUser(user *domain.User) error {
	// Fetch recent submissions
	submissions, err := s.cfClient.GetUserSubmissions(user.CodeforcesHandle, 5000)
	if err != nil {
		return err
	}

	// Fetch updated user info
	userInfo, err := s.cfClient.GetUserInfo(user.CodeforcesHandle)
	if err != nil {
		log.Printf("Warning: Could not fetch user info for %s: %v", user.CodeforcesHandle, err)
	} else {
		user.Rating = userInfo.Rating
		user.Rank = userInfo.Rank
	}

	// Store new submissions
	if err := s.storeSubmissions(user.ID, submissions); err != nil {
		return err
	}

	// Calculate and update streak
	streak := s.calculateStreak(submissions)
	now := time.Now()

	user.CurrentStreak = streak
	if streak > user.MaxStreak {
		user.MaxStreak = streak
	}

	if len(submissions) > 0 {
		latestSubmission := submissions[0]
		submissionTime := time.Unix(latestSubmission.CreationTimeSeconds, 0)
		user.LastSubmissionAt = &submissionTime
	}

	user.LastCheckedAt = &now
	user.TotalSubmissions = len(submissions)

	return s.userRepo.Update(user)
}

func (s *syncService) storeSubmissions(userID uint, cfSubmissions []domain.CodeforcesSubmission) error {
	var newSubmissions []domain.Submission

	for _, cfSub := range cfSubmissions {
		// Check if submission already exists
		_, err := s.submissionRepo.FindByCodeforcesID(int64(cfSub.ID))
		if err == nil {
			continue // Already exists
		}

		submission := domain.Submission{
			UserID:                 userID,
			CodeforcesSubmissionID: int64(cfSub.ID),
			//ProblemName:            cfSub.Problem.Name,
			//ContestID:              cfSub.ContestID,
			//ProblemIndex:           cfSub.Problem.Index,
			Verdict:                cfSub.Verdict,
			//ProgrammingLanguage:    cfSub.ProgrammingLanguage,
			SubmittedAt:            time.Unix(cfSub.CreationTimeSeconds, 0),
		}

		newSubmissions = append(newSubmissions, submission)
	}

	if len(newSubmissions) > 0 {
		return s.submissionRepo.BulkCreate(newSubmissions)
	}

	return nil
}

func (s *syncService) calculateStreak(submissions []domain.CodeforcesSubmission) int {
	if len(submissions) == 0 {
		return 0
	}

	tehranLoc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		panic(err)
	}

	submissionsByDate := make(map[string]bool)
	for _, sub := range submissions {
		if sub.Verdict == "OK" {
			t := time.Unix(sub.CreationTimeSeconds, 0)
			localDate := t.In(tehranLoc).Format("2006-01-02")
			submissionsByDate[localDate] = true
		}
	}

	now := time.Now().In(tehranLoc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tehranLoc)
	yesterday := today.AddDate(0, 0, -1)

	streak := 0
	currentDate := yesterday

	for {
		dateStr := currentDate.Format("2006-01-02")

		if submissionsByDate[dateStr] {
			streak++
			currentDate = currentDate.AddDate(0, 0, -1)
		} else {
			// Gap found → stop counting backward
			break
		}
	}

	// Finally, check if today has a submission — if yes, extend the streak
	if submissionsByDate[today.Format("2006-01-02")] {
		streak++
	}

	return streak
}
