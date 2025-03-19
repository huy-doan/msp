package entities

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user entity in the system
type User struct {
	ID            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Email         string    `json:"email" gorm:"type:varchar(255);uniqueIndex"`
	PasswordHash  string    `json:"-" gorm:"column:password_hash;type:varchar(255)"` // Never exposed in JSON
	RoleID        int       `json:"role_id" gorm:"type:int;not null"`
	Role          *Role     `json:"role" gorm:"foreignKey:RoleID"`
	EnabledMFA    bool      `json:"enabled_mfa" gorm:"type:tinyint(1);default:1"`
	MFATypeID     *int      `json:"mfa_type_id" gorm:"type:int"`
	MFAType       *MFAType  `json:"mfa_type" gorm:"foreignKey:MFATypeID"`
	LastName      string    `json:"last_name" gorm:"type:varchar(100);not null"`
	FirstName     string    `json:"first_name" gorm:"type:varchar(100);not null"`
	LastNameKana  string    `json:"last_name_kana" gorm:"type:varchar(100);not null"`
	FirstNameKana string    `json:"first_name_kana" gorm:"type:varchar(100);not null"`
	AvatarURL     *string   `json:"avatar_url,omitempty" gorm:"type:varchar(255)"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName specifies the database table name
func (User) TableName() string {
	return "users"
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// FullNameKana returns the user's full name in kana
func (u *User) FullNameKana() string {
	return u.FirstNameKana + " " + u.LastNameKana
}

// NewUser creates a new user with the given details
func NewUser(email, password, firstName, lastName, firstNameKana, lastNameKana string, roleID int) (*User, error) {
	// Basic validation
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}
	if password == "" {
		return nil, errors.New("password cannot be empty")
	}
	if firstName == "" || lastName == "" {
		return nil, errors.New("first name and last name cannot be empty")
	}
	if firstNameKana == "" || lastNameKana == "" {
		return nil, errors.New("first name kana and last name kana cannot be empty")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		Email:          email,
		PasswordHash:   string(hashedPassword),
		RoleID:         roleID,
		EnabledMFA:     true, // Default to enabled
		FirstName:      firstName,
		LastName:       lastName,
		FirstNameKana:  firstNameKana,
		LastNameKana:   lastNameKana,
	}, nil
}

// VerifyPassword verifies the provided password against the stored hash
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
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

	u.PasswordHash = string(hashedPassword)
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateProfile updates the user's profile information
func (u *User) UpdateProfile(firstName, lastName, firstNameKana, lastNameKana string) error {
	if firstName == "" || lastName == "" {
		return errors.New("first name and last name cannot be empty")
	}
	
	if firstNameKana == "" || lastNameKana == "" {
		return errors.New("first name kana and last name kana cannot be empty")
	}

	u.FirstName = firstName
	u.LastName = lastName
	u.FirstNameKana = firstNameKana
	u.LastNameKana = lastNameKana
	u.UpdatedAt = time.Now()
	return nil
}

// SetMFA configures the MFA settings for a user
func (u *User) SetMFA(enabled bool, mfaTypeID *int) {
	u.EnabledMFA = enabled
	u.MFATypeID = mfaTypeID
	u.UpdatedAt = time.Now()
}

// IsAdmin checks if the user has admin privileges
func (u *User) IsAdmin() bool {
	return u.Role != nil && u.Role.IsAdmin()
}

// IsNormalUser checks if the user is a customer
func (u *User) IsNormalUser() bool {
	return u.Role != nil && u.Role.IsNormalUser()
}
