package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	tableNmae = "user"
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
	DeletedAt gorm.DeletedAt `gorm:"index;"`
}

// TableName get database table name
func (*UserSchema) TableName() string {
	return tableNmae
}
