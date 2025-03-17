package mysql

import (
	"time"

	"github.com/vnlab/makeshop-payment/src/domain/entities"
)

// UserModel is the MySQL model for User entity
type UserModel struct {
	ID        string    `gorm:"type:varchar(36);primary_key"`
	Username  string    `gorm:"type:varchar(100);uniqueIndex"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex"`
	Password  string    `gorm:"type:varchar(255)"`
	FullName  string    `gorm:"type:varchar(255)"`
	Role      string    `gorm:"type:varchar(50)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName sets the table name for UserModel
func (UserModel) TableName() string {
	return "users"
}

// ToEntity converts UserModel to domain User entity
func (u *UserModel) ToEntity() *entities.User {
	return &entities.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		FullName:  u.FullName,
		Role:      entities.Role(u.Role),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// FromEntity converts domain User entity to UserModel
func (u *UserModel) FromEntity(user *entities.User) {
	u.ID = user.ID
	u.Username = user.Username
	u.Email = user.Email
	u.Password = user.Password
	u.FullName = user.FullName
	u.Role = string(user.Role)
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
}
