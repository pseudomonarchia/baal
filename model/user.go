package model

import (
	"time"

	"gorm.io/gorm"
)

// UserSchema for GORM
type UserSchema struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Username  string `gorm:"size:30;not null"`
	Nickname  string `gorm:"size:30;not null"`
	Email     string `gorm:"size:50;not null"`
	Enable    bool   `gorm:"default:false;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName is GORM hook
func (*UserSchema) TableName() string {
	return "user"
}
