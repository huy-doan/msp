package repositories

import (
	"context"
	"errors"
	"math"

	"github.com/vnlab/makeshop-payment/src/domain/entities"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
	"gorm.io/gorm"
)

// UserRepositoryImpl implements the UserRepository interface
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

// FindByID finds a user by ID
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id int) (*entities.User, error) {
	var user entities.User
	result := r.db.Preload("Role").First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if user not found
		}
		return nil, result.Error
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	result := r.db.Preload("Role").Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if user not found
		}
		return nil, result.Error
	}
	return &user, nil
}

// Create creates a new user
func (r *UserRepositoryImpl) Create(ctx context.Context, user *entities.User) error {
	return r.db.Create(user).Error
}

// Update updates an existing user
func (r *UserRepositoryImpl) Update(ctx context.Context, user *entities.User) error {
	return r.db.Save(user).Error
}

// Delete soft-deletes a user by ID
func (r *UserRepositoryImpl) Delete(ctx context.Context, id int) error {
	return r.db.Delete(&entities.User{}, id).Error
}

// List lists all users with pagination
func (r *UserRepositoryImpl) List(ctx context.Context, page, pageSize int) ([]*entities.User, int, error) {
	var users []*entities.User
	var count int64

	// Count total records
	if err := r.db.Model(&entities.User{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := r.db.Preload("Role").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	totalPages := int(math.Ceil(float64(count) / float64(pageSize)))
	return users, totalPages, nil
}
