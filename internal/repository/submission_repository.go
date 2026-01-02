package repository

import (
	"time"

	"github.com/pouyatavakoli/CodeStreaks-web/internal/domain"
	"gorm.io/gorm"
)

type SubmissionRepository interface {
	Create(submission *domain.Submission) error
	BulkCreate(submissions []domain.Submission) error
	FindByCodeforcesID(cfID int64) (*domain.Submission, error)
	GetUserSubmissions(userID uint, limit int) ([]domain.Submission, error)
	GetLatestSubmissionForUser(userID uint) (*domain.Submission, error)
	GetSubmissionsAfter(userID uint, after time.Time) ([]domain.Submission, error)
}

type submissionRepository struct {
	db *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{db: db}
}

func (r *submissionRepository) Create(submission *domain.Submission) error {
	return r.db.Create(submission).Error
}

func (r *submissionRepository) BulkCreate(submissions []domain.Submission) error {
	if len(submissions) == 0 {
		return nil
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Use CreateInBatches for better performance
		return tx.CreateInBatches(submissions, 100).Error
	})
}

func (r *submissionRepository) FindByCodeforcesID(cfID int64) (*domain.Submission, error) {
	var submission domain.Submission
	err := r.db.Where("codeforces_submission_id = ?", cfID).First(&submission).Error
	if err != nil {
		return nil, err
	}
	return &submission, nil
}

func (r *submissionRepository) GetUserSubmissions(userID uint, limit int) ([]domain.Submission, error) {
	var submissions []domain.Submission
	err := r.db.Where("user_id = ?", userID).
		Order("submitted_at DESC").
		Limit(limit).
		Find(&submissions).Error
	return submissions, err
}

func (r *submissionRepository) GetLatestSubmissionForUser(userID uint) (*domain.Submission, error) {
	var submission domain.Submission
	err := r.db.Where("user_id = ?", userID).
		Order("submitted_at DESC").
		First(&submission).Error
	if err != nil {
		return nil, err
	}
	return &submission, nil
}

func (r *submissionRepository) GetSubmissionsAfter(userID uint, after time.Time) ([]domain.Submission, error) {
	var submissions []domain.Submission
	err := r.db.Where("user_id = ? AND submitted_at > ?", userID, after).
		Order("submitted_at ASC").
		Find(&submissions).Error
	return submissions, err
}
