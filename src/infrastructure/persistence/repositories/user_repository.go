package repositories

import (
	"context"
	"errors"
	"math"

	"github.com/vnlab/makeshop-payment/src/domain/entities"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
	"gorm.io/gorm"
)

// UserRepositoryInfra implements the UserRepository interface
type UserRepositoryInfra struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &UserRepositoryInfra{
		db: db,
	}
}

// FindByID finds a user by ID
func (r *UserRepositoryInfra) FindByID(ctx context.Context, id string) (model *entities.User, err error) {
	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if user not found
		}
		return nil, err
	}
	return model, nil
}

// FindByUsername finds a user by username
func (r *UserRepositoryInfra) FindByUsername(ctx context.Context, username string) (model *entities.User, err error) {
	if err := r.db.First(&model, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if user not found
		}
		return nil, err
	}

	return model, nil
}

// FindByEmail finds a user by email
func (r *UserRepositoryInfra) FindByEmail(ctx context.Context, email string) (model *entities.User, err error) {
	if err := r.db.First(&model, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if user not found
		}
		return nil, err
	}
	return model, nil
}

// Create creates a new user
func (r *UserRepositoryInfra) Create(ctx context.Context, user *entities.User) error {
	model := &entities.User{}
	return r.db.Create(model).Error
}

// Update updates an existing user
func (r *UserRepositoryInfra) Update(ctx context.Context, user *entities.User) error {
	model := &entities.User{}
	return r.db.Save(model).Error
}

// Delete deletes a user by ID
func (r *UserRepositoryInfra) Delete(ctx context.Context, id string) error {
	return r.db.Delete(&entities.User{}, "id = ?", id).Error
}

// List lists all users with pagination
func (r *UserRepositoryInfra) List(ctx context.Context, page, pageSize int) ([]*entities.User, int, error) {
	var models []*entities.User
	var count int64

	// Count total records
	if err := r.db.Model(&entities.User{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := r.db.Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	// Convert to domain entities
	users := make([]*entities.User, len(models))
	for i, model := range models {
		users[i] = model
	}

	totalPages := int(math.Ceil(float64(count) / float64(pageSize)))
	return users, totalPages, nil
}
