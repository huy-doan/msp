package repositories

import (
	"context"

	"github.com/vnlab/makeshop-payment/src/domain/entities"
)

// RoleRepository defines the interface for role data access
type RoleRepository interface {
	// FindByID finds a role by ID
	FindByID(ctx context.Context, id int) (*entities.Role, error)

	// FindByCode finds a role by code
	FindByCode(ctx context.Context, code string) (*entities.Role, error)
}
