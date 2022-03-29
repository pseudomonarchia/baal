package model

import "time"

// UserHashSchema for GORM
type UserHashSchema struct {
	UserID    uint   `gorm:"primaryKey"`
	Hash      string `gorm:"size:100"`
	UpdatedAt time.Time
}

// TableName get database table name
func (*UserHashSchema) TableName() string {
	return "user_hash"
}
