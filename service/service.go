package service

import "gorm.io/gorm"

// Services represents a global Service struct
type Services struct {
	OAuth OAuthFace
	User  UserFace
}

// New return all service
func New(db *gorm.DB) *Services {
	return &Services{
		OAuth: &OAuth{db},
		User:  &User{db},
	}
}
