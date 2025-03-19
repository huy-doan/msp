package entities

import (
	"time"
)

// MFAType represents a multi-factor authentication type
type MFAType struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	No        int       `json:"no" gorm:"type:int;not null"`
	Title     string    `json:"title" gorm:"type:varchar(255);not null"`
	IsActive  int       `json:"is_active" gorm:"type:int;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName specifies the database table name
func (MFAType) TableName() string {
	return "master_mfa_types"
}

// IsActiveType checks if this MFA type is active
func (m *MFAType) IsActiveType() bool {
	return m.IsActive == 1
}
