package resolvers

import (
	"context"
	"errors"

	"github.com/vnlab/makeshop-payment/src/api/graphql/generated"
	"github.com/vnlab/makeshop-payment/src/api/graphql/middleware"
	"github.com/vnlab/makeshop-payment/src/domain/models"
)

// Query returns the QueryResolver implementation
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{Resolver: r}
}

type queryResolver struct {
	*Resolver
}

// Define custom errors
var (
	ErrNotAuthenticated = errors.New("not authenticated")
	ErrForbidden        = errors.New("forbidden")
)

// Me returns the currently authenticated user
func (r *queryResolver) Me(ctx context.Context) (*models.User, error) {
	// Check auth
	err := middleware.CheckAuth(ctx)
	if err != nil {
		return nil, ErrNotAuthenticated
	}
	
	userId, err := middleware.GetUserID(ctx)
	if err != nil {
		return nil, ErrNotAuthenticated
	}

	user, err := r.userUsecase.GetUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// User returns a user by ID
func (r *queryResolver) User(ctx context.Context, id int) (*models.User, error) {
	// Check auth
	err := middleware.CheckAuth(ctx)
	if err != nil {
		return nil, ErrNotAuthenticated
	}
	
	// Check if the user has admin rights
	if !middleware.IsAdminRole(ctx) {
		return nil, ErrForbidden
	}

	user, err := r.userUsecase.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Users returns a paginated list of users
func (r *queryResolver) Users(ctx context.Context, page *int, pageSize *int) (*generated.PaginatedUsers, error) {
	// Check auth
	err := middleware.CheckAuth(ctx)
	if err != nil {
		return nil, ErrNotAuthenticated
	}
	
	// Check if the user has admin rights
	if !middleware.IsAdminRole(ctx) {
		return nil, ErrForbidden
	}

	p := 1
	if page != nil {
		p = *page
	}

	ps := 10
	if pageSize != nil {
		ps = *pageSize
	}

	users, totalPages, err := r.userUsecase.ListUsers(ctx, p, ps)
	if err != nil {
		return nil, err
	}

	return &generated.PaginatedUsers{
		Users:      users,
		Page:       p,
		PageSize:   ps,
		TotalPages: totalPages,
	}, nil
}
