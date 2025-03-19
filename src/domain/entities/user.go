package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Role represents user roles
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleCustomer Role = "customer"
)

// User represents a user entity
type User struct {
	ID        string    `json:"id" gorm:"type:varchar(36);primary_key"`
	Username  string    `json:"username" gorm:"type:varchar(100);uniqueIndex"`
	Email     string    `json:"email" gorm:"type:varchar(255);uniqueIndex"`
	Password  string    `json:"-" gorm:"type:varchar(255)"` // Password is never exposed in JSON
	FullName  string    `json:"full_name" gorm:"type:varchar(255)"`
	Role      Role      `json:"role" gorm:"type:varchar(50)"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the database table name
func (User) TableName() string {
	return "users"
}

// NewUser creates a new user with the given details
func NewUser(username, email, password, fullName string, role Role) (*User, error) {
	// Validation logic unchanged
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	if password == "" {
		return nil, errors.New("password cannot be empty")
	}

	if fullName == "" {
		return nil, errors.New("full name cannot be empty")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		ID:        uuid.New().String(),
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		FullName:  fullName,
		Role:      role,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// VerifyPassword verifies the provided password against the stored hash
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// ChangePassword changes the user's password
func (u *User) ChangePassword(newPassword string) error {
	if newPassword == "" {
		return errors.New("new password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateProfile updates the user's profile information
func (u *User) UpdateProfile(fullName string) error {
	if fullName == "" {
		return errors.New("full name cannot be empty")
	}

	u.FullName = fullName
	u.UpdatedAt = time.Now()
	return nil
}
