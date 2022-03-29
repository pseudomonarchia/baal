package model

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	// OAuthProviderGoogle ...
	OAuthProviderGoogle = "google"
)

// OAuthSchema for GORM
type OAuthSchema struct {
	UID       string `gorm:"primaryKey"`
	UserID    uint
	Provider  string `gorm:"size:36"`
	TokenInfo datatypes.JSON
}

// GoogleOAuthRequest ...
type GoogleOAuthRequest struct {
	Redirect string `form:"redirect" validate:"required,url"`
}

// GoogleOAuthResponse ...
type GoogleOAuthResponse struct {
	State    string `form:"state" validate:"required"`
	Code     string `form:"code" validate:"required"`
	Scope    string `form:"scope" validate:"required"`
	AuthUser string `form:"authuser" validate:"required"`
	Prompt   string `form:"prompt" validate:"required"`
}

// GoogleOAuthUserInfo ...
type GoogleOAuthUserInfo struct {
	Sub           string `form:"sub"`
	Name          string `form:"name"`
	FirstName     string `form:"given_name"`
	LastName      string `form:"family_name"`
	Picture       string `form:"picture"`
	Email         string `form:"email"`
	EmailVerified bool   `form:"email_verified"`
	Locale        string `form:"Locale"`
}

// TableName is GORM hook
func (*OAuthSchema) TableName() string {
	return "oauth"
}

// BeforeCreate is GORM hook
func (o *OAuthSchema) BeforeCreate(tx *gorm.DB) error {
	o.UID = uuid.NewString()
	return nil
}
