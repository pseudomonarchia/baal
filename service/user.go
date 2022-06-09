package service

import (
	"baal/model"

	"gorm.io/gorm"
)

// UserFace User service interface
type UserFace interface {
	GetByQuery(m *model.UserSchema) (*model.UserSchema, error)
}

// User ...
type User struct {
	Database *gorm.DB
}

// GetByQuery fetch user data by query
func (u *User) GetByQuery(m *model.UserSchema) (*model.UserSchema, error) {
	user := &model.UserSchema{}
	err := u.Database.Take(user).Error
	return user, err
}
