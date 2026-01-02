package repository

import (
	"time"

	"github.com/pouyatavakoli/CodeStreaks-web/internal/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *domain.User) error
	Update(user *domain.User) error
	FindByHandle(handle string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
	GetLeaderboard(limit, offset int) ([]domain.User, error)
	GetAllActiveUsers() ([]domain.User, error)
	CountUsers() (int64, error)
	BulkUpdate(users []domain.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) FindByHandle(handle string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("codeforces_handle = ?", handle).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetLeaderboard(limit, offset int) ([]domain.User, error) {
	var users []domain.User
	err := r.db.Where("is_active = ?", true).
		Order("current_streak DESC, max_streak DESC, rating DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error
	return users, err
}

func (r *userRepository) GetAllActiveUsers() ([]domain.User, error) {
	var users []domain.User
	err := r.db.Where("is_active = ?", true).Find(&users).Error
	return users, err
}

func (r *userRepository) CountUsers() (int64, error) {
	var count int64
	err := r.db.Model(&domain.User{}).Where("is_active = ?", true).Count(&count).Error
	return count, err
}

func (r *userRepository) BulkUpdate(users []domain.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, user := range users {
			user.UpdatedAt = time.Now()
			if err := tx.Save(&user).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
