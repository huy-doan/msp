package repositories

import (
	"context"

	"github.com/vnlab/makeshop-payment/src/domain/entities"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// FindByID finds a user by ID
	FindByID(ctx context.Context, id string) (*entities.User, error)

	// FindByUsername finds a user by username
	FindByUsername(ctx context.Context, username string) (*entities.User, error)

	// FindByEmail finds a user by email
	FindByEmail(ctx context.Context, email string) (*entities.User, error)

	// Create creates a new user
	Create(ctx context.Context, user *entities.User) error

	// Update updates an existing user
	Update(ctx context.Context, user *entities.User) error

	// Delete deletes a user by ID
	Delete(ctx context.Context, id string) error

	// List lists all users with pagination
	List(ctx context.Context, page, pageSize int) ([]*entities.User, int, error)
}
