package resolvers

import (
	"context"

	"github.com/vnlab/makeshop-payment/src/api/graphql/generated"
	"github.com/vnlab/makeshop-payment/src/api/graphql/middleware"
	"github.com/vnlab/makeshop-payment/src/domain/entities"
	"github.com/vnlab/makeshop-payment/src/usecase"
)

// Mutation returns MutationResolver implementation
func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{Resolver: r}
}

type mutationResolver struct {
	*Resolver
}

// Login implements the login mutation
func (r *mutationResolver) Login(ctx context.Context, input generated.LoginInput) (*generated.AuthResponse, error) {
	loginReq := usecase.LoginRequest{
		Username: input.Username,
		Password: input.Password,
	}

	loginResp, err := r.userUsecase.Login(ctx, loginReq)
	if err != nil {
		return nil, err
	}

	return &generated.AuthResponse{
		Token: loginResp.Token,
		User:  loginResp.User,
	}, nil
}

// Register implements the register mutation
func (r *mutationResolver) Register(ctx context.Context, input generated.RegisterInput) (*entities.User, error) {
	registerReq := usecase.RegisterRequest{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
		FullName: input.FullName,
	}

	user, err := r.userUsecase.Register(ctx, registerReq)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateProfile implements the updateProfile mutation
func (r *mutationResolver) UpdateProfile(ctx context.Context, input generated.UpdateProfileInput) (*entities.User, error) {
	err := middleware.CheckAuth(ctx)
	if err != nil {
		return nil, ErrNotAuthenticated
	}
	userId, err := middleware.GetUserID(ctx)
	if err != nil {
		return nil, ErrNotAuthenticated
	}

	user, err := r.userUsecase.UpdateUserProfile(ctx, userId, input.FullName)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword implements the changePassword mutation
func (r *mutationResolver) ChangePassword(ctx context.Context, input generated.ChangePasswordInput) (bool, error) {
	err := middleware.CheckAuth(ctx)
	if err != nil {
		return false, ErrNotAuthenticated
	}
	userId, err := middleware.GetUserID(ctx)
	if err != nil {
		return false, ErrNotAuthenticated
	}

	err = r.userUsecase.ChangePassword(ctx, userId, input.CurrentPassword, input.NewPassword)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Logout implements the logout mutation
func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	err := middleware.CheckAuth(ctx)
	if err != nil {
		return false, ErrNotAuthenticated
	}

	// get token from context (added token to context in middleware)
	token, ok := ctx.Value("token").(string)
	if ok && token != "" {
		// Thêm token vào blacklist
		r.jwtService.BlacklistToken(token)
	}
	
	// return true to confirm successful logout
	return true, nil
}
