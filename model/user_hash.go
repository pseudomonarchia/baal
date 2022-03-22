package model

import "time"

const (
	tableName = "user_hash"
)

// UserHashSchema for GORM
type UserHashSchema struct {
	UserID    uint   `gorm:"primaryKey"`
	Hash      string `gorm:"size:100"`
	UpdatedAt time.Time
}

// TableName get database table name
func (*UserHashSchema) TableName() string {
	return tableName
}
