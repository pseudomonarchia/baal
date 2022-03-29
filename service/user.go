package service

import (
	"baal/model"
	"errors"

	"gorm.io/gorm"
)

// UserFace User service interface
type UserFace interface {
	GetByQuery(m *model.UserSchema) (*model.UserSchema, bool)
}

// User ...
type User struct {
	Datebase *gorm.DB
}

// GetByQuery fetch user data by query
func (u *User) GetByQuery(m *model.UserSchema) (*model.UserSchema, bool) {
	user := &model.UserSchema{}
	r := u.Datebase.First(user, m)
	return user, errors.Is(r.Error, gorm.ErrRecordNotFound)
}
