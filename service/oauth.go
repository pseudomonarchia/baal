package service

import (
	"baal/config"
	"baal/model"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	callbackURL  = "/api/v1/oauth/callback"
	openIDScope  = "openid"
	emailScope   = "https://www.googleapis.com/auth/userinfo.email"
	profileScope = "https://www.googleapis.com/auth/userinfo.profile"
	userinfoURL  = "https://www.googleapis.com/oauth2/v3/userinfo"
)

// OAuthFace OAuth service interface
type OAuthFace interface {
	NewState() string
	GetLoginURL(requestURL, state string) string
	GetToken(requestURL, code string) (*oauth2.Token, error)
	GetInfo(requestURL string, token *oauth2.Token) (*model.GoogleOAuthUserInfo, error)
	RefreshToken(requestURL string, token *oauth2.Token) (*oauth2.Token, error)
	SaveToken(userID uint, UID string, token *oauth2.Token) (*model.OAuthTokenSchema, error)
	GenerateUID() string
	GenerateHashFromUID(UID string) string
	GenerateRefreshToken(OAuthUser *model.OAuthTokenSchema, IP string) (*model.OAuthRefreshSchema, error)
	FindToken(UID string) *model.OAuthTokenSchema
	FindRefreshTokenFormUID(token string, UID string) (*model.OAuthRefreshSchema, error)
	CleanToken(UID string) error
	UseToken(tx *gorm.DB, token model.OAuthTokenSchema) error
	HashTokenToSHA(r *datatypes.JSON) []byte
	CheckTokenAndReplace(authorizationHeader string) (string, bool)
	DecodeToken(str string) (*jwt.StandardClaims, error)
	ValidateToken(str string, r *datatypes.JSON) (bool, error)
	SignToken(payload *jwt.StandardClaims, r *datatypes.JSON) string
}

// OAuth ...
type OAuth struct {
	Database *gorm.DB
}

func (*OAuth) googleOAuth(requestURL string) *oauth2.Config {
	redirectURL := fmt.Sprintf(
		"%s://%s%s",
		config.Global.PROTOCOL(),
		requestURL,
		callbackURL,
	)

	return &oauth2.Config{
		ClientID:     config.Secret.Oauth.Google.ClientID,
		ClientSecret: config.Secret.Oauth.Google.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  redirectURL,
		Scopes: []string{
			openIDScope,
			emailScope,
			profileScope,
		},
	}
}

// NewState generate a new state for oauth
func (*OAuth) NewState() string {
	b := make([]byte, 16)
	rand.Read(b)

	return base64.URLEncoding.EncodeToString(b)
}

// GetLoginURL generate a login url
func (o *OAuth) GetLoginURL(requestURL, state string) string {
	return o.googleOAuth(requestURL).AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// GetToken bring in code to get token
func (o *OAuth) GetToken(requestURL, code string) (*oauth2.Token, error) {
	return o.googleOAuth(requestURL).Exchange(oauth2.NoContext, code)
}

// GetInfo carry token to get social user information
func (o *OAuth) GetInfo(
	requestURL string,
	token *oauth2.Token,
) (
	*model.GoogleOAuthUserInfo,
	error,
) {
	client := o.googleOAuth(requestURL).Client(oauth2.NoContext, token)
	res, _ := client.Get(userinfoURL)
	defer res.Body.Close()

	user := &model.GoogleOAuthUserInfo{}
	err := json.NewDecoder(res.Body).Decode(user)

	return user, err
}

// RefreshToken ...
func (o *OAuth) RefreshToken(
	requestURL string,
	token *oauth2.Token,
) (
	*oauth2.Token,
	error,
) {
	return o.
		googleOAuth(requestURL).
		TokenSource(oauth2.NoContext, token).
		Token()
}

// SaveToken store the token in database and delete the existing token of this user
func (o *OAuth) SaveToken(
	userID uint,
	UID string,
	token *oauth2.Token,
) (
	*model.OAuthTokenSchema,
	error,
) {
	b, _ := json.Marshal(token)
	OAuthUser := &model.OAuthTokenSchema{
		UID:       UID,
		UserID:    userID,
		TokenInfo: datatypes.JSON(b),
		Provider:  model.OAuthProviderGoogle,
	}

	err := o.Database.Save(OAuthUser).Error
	return OAuthUser, err
}

// GenerateUID ...
func (o *OAuth) GenerateUID() string {
	return uuid.NewString()
}

// GenerateHashFromUID ..
func (o *OAuth) GenerateHashFromUID(UID string) string {
	b, _ := bcrypt.GenerateFromPassword([]byte(UID), bcrypt.DefaultCost)
	return string(b)
}

// GenerateRefreshToken ...
func (o *OAuth) GenerateRefreshToken(
	OAuthUser *model.OAuthTokenSchema,
	IP string,
) (
	*model.OAuthRefreshSchema,
	error,
) {
	refreshInfo := &model.OAuthRefreshSchema{
		OAuthUID:  OAuthUser.UID,
		IP:        IP,
		Token:     o.GenerateHashFromUID(OAuthUser.UID),
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(24 * 7 * time.Hour),
	}

	err := o.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(refreshInfo).Error; err != nil {
			return err
		}

		if err := o.UseToken(tx, *OAuthUser); err != nil {
			return err
		}

		return nil
	})

	return refreshInfo, err
}

// FindToken ...
func (o *OAuth) FindToken(UID string) *model.OAuthTokenSchema {
	OAuthUser := &model.OAuthTokenSchema{UID: UID, Use: false}
	record := o.Database.Where(OAuthUser, "Use").Take(OAuthUser)
	if record.Error == gorm.ErrRecordNotFound {
		return nil
	}

	return OAuthUser
}

// FindRefreshTokenFormUID ...
func (o *OAuth) FindRefreshTokenFormUID(
	token string,
	UID string,
) (
	*model.OAuthRefreshSchema,
	error,
) {
	refreshToken := &model.OAuthRefreshSchema{Token: token, OAuthUID: UID}
	err := o.Database.
		Preload("OAuthToken").
		Preload("OAuthToken.User").
		Where(refreshToken, "Token", "UID").
		Take(refreshToken).
		Error

	return refreshToken, err
}

// CleanToken ...
func (o *OAuth) CleanToken(UID string) error {
	OAuthUser := &model.OAuthTokenSchema{UID: UID}
	err := o.Database.Delete(OAuthUser).Error

	return err
}

// UseToken ...
func (*OAuth) UseToken(tx *gorm.DB, token model.OAuthTokenSchema) error {
	token.Use = true
	err := tx.Save(token).Error

	return err
}

// HashTokenToSHA ...
func (*OAuth) HashTokenToSHA(r *datatypes.JSON) []byte {
	h := sha1.New()
	h.Write([]byte(r.String()))
	secret := h.Sum(nil)

	return secret
}

// CheckTokenAndReplace ...
func (*OAuth) CheckTokenAndReplace(bearerToken string) (string, bool) {
	reg := regexp.MustCompile(`(^Bearer )(.*\..*\..*$)`)
	validate := reg.Match([]byte(bearerToken))
	str := reg.ReplaceAllString(bearerToken, "$2")

	return str, validate
}

// DecodeToken ...
func (*OAuth) DecodeToken(str string) (*jwt.StandardClaims, error) {
	token := &jwt.StandardClaims{}
	s := strings.Split(str, ".")[1]
	b, err := jwt.DecodeSegment(s)
	if err != nil {
		return token, err
	}

	err = json.Unmarshal(b, token)
	return token, err
}

// ValidateToken ...
func (o *OAuth) ValidateToken(str string, r *datatypes.JSON) (bool, error) {
	var keyFunc = func(t *jwt.Token) (interface{}, error) {
		return o.HashTokenToSHA(r), nil
	}

	_, err := jwt.Parse(str, keyFunc)
	return err == nil, err
}

// SignToken ...
func (o *OAuth) SignToken(payload *jwt.StandardClaims, r *datatypes.JSON) string {
	secret := o.HashTokenToSHA(r)
	JWT, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString(secret)
	return JWT
}
