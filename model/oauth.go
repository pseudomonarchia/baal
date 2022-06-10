package model

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
	"gorm.io/datatypes"
)

// OAuthProvider ...
type OAuthProvider string

// GrantType ...
type GrantType string

// const ...
const (
	OAuthTokenTable                         = "oauth_token"
	OAuthRefreshTable                       = "oauth_refresh"
	OAuthProviderGoogle       OAuthProvider = "google"
	GrantTypeFromCode         GrantType     = "code"
	GrantTypeFromRefreshToken GrantType     = "refresh_token"
)

// OAuthTokenSchema for GORM
type OAuthTokenSchema struct {
	UID       string         `gorm:"primaryKey"`
	UserID    uint           `gorm:"not null"`
	Provider  OAuthProvider  `gorm:"size:10"`
	Use       bool           `gorm:"default:false"`
	TokenInfo datatypes.JSON `gorm:"not null"`
	User      UserSchema     `gorm:"foreignKey:UserID;references:ID"`
}

// OAuthRefreshSchema for GORM
type OAuthRefreshSchema struct {
	OAuthUID   string           `gorm:"primaryKey;not null;column:oauth_uid"`
	IP         string           `gorm:"size:30;not null"`
	Token      string           `gorm:"unique"`
	IssuedAt   time.Time        `gorm:"not null"`
	ExpiresAt  time.Time        `gorm:"not null"`
	OAuthToken OAuthTokenSchema `gorm:"foreignKey:OAuthUID;references:UID"`
}

// GoogleOAuthRequest ...
type GoogleOAuthRequest struct {
	Redirect string `form:"redirect" validate:"required,url"`
}

// TokenRequest ...
type TokenRequest struct {
	GrantType GrantType `json:"grant_type" form:"grant_type" validate:"required,oneof=code refresh_token"`
	Code      string    `form:"code" validate:"required"`
}

// TokenSchema ...
type TokenSchema struct {
	Expiry       time.Time `json:"expiry"`
	TokenType    string    `json:"token_type"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
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

// JWTClaims ...
type JWTClaims struct {
	jwt.StandardClaims
}

// TableName is GORM hook
func (*OAuthTokenSchema) TableName() string {
	return OAuthTokenTable
}

// UnmarshalToken ...
func (o *OAuthTokenSchema) UnmarshalToken() (*oauth2.Token, error) {
	token := &oauth2.Token{}
	err := json.Unmarshal([]byte(o.TokenInfo.String()), token)

	return token, err
}

// TableName is GORM hook
func (*OAuthRefreshSchema) TableName() string {
	return OAuthRefreshTable
}

// SignToken ...
func (j *JWTClaims) SignToken(secret interface{}) string {
	JWT, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, j).SignedString(secret)
	return JWT
}

// DecodeToken ...
func (j *JWTClaims) DecodeToken(token string) error {
	s := strings.Split(token, ".")[1]
	b, err := jwt.DecodeSegment(s)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, j)
	return err
}

// CheckBearerToken ...
func (*JWTClaims) CheckBearerToken(bearerToken string) bool {
	reg := regexp.MustCompile(`(^Bearer )(.*\..*\..*$)`)
	validate := reg.Match([]byte(bearerToken))

	return validate
}

// PickTokenFormBearer ...
func (*JWTClaims) PickTokenFormBearer(bearerToken string) string {
	reg := regexp.MustCompile(`(^Bearer )(.*\..*\..*$)`)
	str := reg.ReplaceAllString(bearerToken, "$2")

	return str
}

// ValidateToken ...
func (j *JWTClaims) ValidateToken(token string, secret interface{}) error {
	str := j.PickTokenFormBearer(token)
	var keyFunc = func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	}

	_, err := jwt.Parse(str, keyFunc)
	return err
}
