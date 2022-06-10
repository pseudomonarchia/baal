package service

import (
	"baal/config"
	"baal/model"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

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
	SaveToken(ref *model.OAuthTokenSchema) error
	GenerateUID() string
	GenerateHashFromUID(UID string) string
	GenerateRefreshToken(ref *model.OAuthTokenSchema, IP string) (*model.OAuthRefreshSchema, error)
	FindToken(ref *model.OAuthTokenSchema, query ...interface{}) error
	FindRefreshToken(ref *model.OAuthRefreshSchema, query ...interface{}) error
	CleanTokenResource(ref *model.OAuthRefreshSchema) error
	HashJSONToSHA(r *datatypes.JSON) []byte
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
func (o *OAuth) SaveToken(ref *model.OAuthTokenSchema) error {
	return o.Database.Save(ref).Error
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
	ref *model.OAuthTokenSchema,
	IP string,
) (
	*model.OAuthRefreshSchema,
	error,
) {
	refreshInfo := &model.OAuthRefreshSchema{
		OAuthUID:   ref.UID,
		IP:         IP,
		Token:      o.GenerateHashFromUID(ref.UID),
		IssuedAt:   time.Now(),
		ExpiresAt:  time.Now().Add(24 * 7 * time.Hour),
		OAuthToken: *ref,
	}

	err := o.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(refreshInfo).Error; err != nil {
			return err
		}

		if ref.Use {
			return nil
		}

		ref.Use = true
		if err := tx.Save(ref).Error; err != nil {
			return err
		}

		return nil
	})

	return refreshInfo, err
}

// FindToken ...
func (o *OAuth) FindToken(ref *model.OAuthTokenSchema, query ...interface{}) error {
	err := o.Database.Where(ref, query...).Take(ref).Error
	return err
}

// FindRefreshToken ...
func (o *OAuth) FindRefreshToken(ref *model.OAuthRefreshSchema, query ...interface{}) error {
	err := o.Database.
		Preload("OAuthToken").
		Preload("OAuthToken.User").
		Where(ref, query...).
		Take(ref).
		Error

	return err
}

// CleanTokenResource ...
func (o *OAuth) CleanTokenResource(ref *model.OAuthRefreshSchema) error {
	err := o.Database.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(ref.OAuthToken).Error
		if err != nil {
			return err
		}

		err = tx.Delete(ref).Error
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// HashJSONToSHA ...
func (*OAuth) HashJSONToSHA(r *datatypes.JSON) []byte {
	h := sha1.New()
	h.Write([]byte(r.String()))
	secret := h.Sum(nil)

	return secret
}
