package service

import (
	"baal/model"

	"gorm.io/gorm"
)

// UserFace User service interface
type UserFace interface {
	GetByQuery(m *model.UserSchema) error
}

// User ...
type User struct {
	Database *gorm.DB
}

// GetByQuery fetch user data by query
func (u *User) GetByQuery(m *model.UserSchema) error {
	err := u.Database.Where(m).Take(m).Error
	return err
}
