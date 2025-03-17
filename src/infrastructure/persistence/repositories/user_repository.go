package repositories

import (
	"context"
	"errors"
	"math"

	"github.com/vnlab/makeshop-payment/src/domain/entities"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
	"github.com/vnlab/makeshop-payment/src/infrastructure/persistence/mysql"
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
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.User, error) {
	var model mysql.UserModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if user not found
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

// FindByUsername finds a user by username
func (r *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (*entities.User, error) {
	var model mysql.UserModel
	if err := r.db.WithContext(ctx).First(&model, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if user not found
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

// FindByEmail finds a user by email
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var model mysql.UserModel
	if err := r.db.WithContext(ctx).First(&model, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if user not found
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

// Create creates a new user
func (r *UserRepositoryImpl) Create(ctx context.Context, user *entities.User) error {
	model := &mysql.UserModel{}
	model.FromEntity(user)
	return r.db.WithContext(ctx).Create(model).Error
}

// Update updates an existing user
func (r *UserRepositoryImpl) Update(ctx context.Context, user *entities.User) error {
	model := &mysql.UserModel{}
	model.FromEntity(user)
	return r.db.WithContext(ctx).Save(model).Error
}

// Delete deletes a user by ID
func (r *UserRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&mysql.UserModel{}, "id = ?", id).Error
}

// List lists all users with pagination
func (r *UserRepositoryImpl) List(ctx context.Context, page, pageSize int) ([]*entities.User, int, error) {
	var models []mysql.UserModel
	var count int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&mysql.UserModel{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	// Convert to domain entities
	users := make([]*entities.User, len(models))
	for i, model := range models {
		users[i] = model.ToEntity()
	}

	totalPages := int(math.Ceil(float64(count) / float64(pageSize)))
	return users, totalPages, nil
}
