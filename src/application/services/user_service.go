package services

import (
	"context"
	"errors"

	"github.com/vnlab/makeshop-payment/src/domain/entities"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
	"github.com/vnlab/makeshop-payment/src/infrastructure/auth"
)

// UserService handles user-related business logic
type UserService struct {
	userRepo   repositories.UserRepository
	jwtService *auth.JWTService
}

// NewUserService creates a new UserService
func NewUserService(userRepo repositories.UserRepository, jwtService *auth.JWTService) *UserService {
	return &UserService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
}

// LoginResponse represents a login response with token
type LoginResponse struct {
	Token string         `json:"token"`
	User  *entities.User `json:"user"`
}

// Login authenticates a user and returns a JWT token
func (s *UserService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid username or password")
	}

	if !user.VerifyPassword(req.Password) {
		return nil, errors.New("invalid username or password")
	}

	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// Register creates a new user
func (s *UserService) Register(ctx context.Context, req RegisterRequest) (*entities.User, error) {
	// Check if username already exists
	existingUser, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	existingUser, err = s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Create new user with customer role by default
	user, err := entities.NewUser(req.Username, req.Email, req.Password, req.FullName, entities.RoleCustomer)
	if err != nil {
		return nil, err
	}

	// Save user to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id string) (*entities.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

// UpdateUserProfile updates a user's profile
func (s *UserService) UpdateUserProfile(ctx context.Context, userID string, fullName string) (*entities.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if err := user.UpdateProfile(fullName); err != nil {
		return nil, err
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	if !user.VerifyPassword(currentPassword) {
		return errors.New("current password is incorrect")
	}

	if err := user.ChangePassword(newPassword); err != nil {
		return err
	}

	return s.userRepo.Update(ctx, user)
}

// ListUsers lists users with pagination
func (s *UserService) ListUsers(ctx context.Context, page, pageSize int) ([]*entities.User, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.userRepo.List(ctx, page, pageSize)
}
