package service

import (
	"baal/config"
	"baal/model"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// OAuthFace OAuth service interface
type OAuthFace interface {
	NewState() string
	GetLoginURL(state string) string
	GetToken(code string) (*oauth2.Token, error)
	GetInfo(token *oauth2.Token) (*model.GoogleOAuthUserInfo, error)
	SaveToken(userID uint, token *oauth2.Token) (*model.OAuthSchema, error)
}

// OAuth ...
type OAuth struct {
	Datebase *gorm.DB
}

func (o *OAuth) googleOAuth() *oauth2.Config {
	redirectURL := fmt.Sprintf("%s/api/v1/login/callback", config.Global.URL())
	return &oauth2.Config{
		ClientID:     "176852869927-31dtie98t8fj0fsmdc7g9em1o1mrkh95.apps.googleusercontent.com",
		ClientSecret: "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",
		Endpoint:     google.Endpoint,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"openid",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
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
func (o *OAuth) GetLoginURL(state string) string {
	return o.googleOAuth().AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// GetToken bring in code to get token
func (o *OAuth) GetToken(code string) (*oauth2.Token, error) {
	return o.googleOAuth().Exchange(oauth2.NoContext, code)
}

// GetInfo carry token to get social user information
func (o *OAuth) GetInfo(token *oauth2.Token) (*model.GoogleOAuthUserInfo, error) {
	client := o.googleOAuth().Client(oauth2.NoContext, token)
	res, _ := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	defer res.Body.Close()

	user := &model.GoogleOAuthUserInfo{}
	err := json.NewDecoder(res.Body).Decode(user)

	return user, err
}

// SaveToken store the token in database and delete the existing token of this user
func (o *OAuth) SaveToken(userID uint, token *oauth2.Token) (*model.OAuthSchema, error) {
	b, _ := json.Marshal(token)
	OAuthUser := &model.OAuthSchema{
		UserID:    userID,
		TokenInfo: datatypes.JSON(b),
		Provider:  model.OAuthProviderGoogle,
	}

	err := o.Datebase.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&model.OAuthSchema{}).Error; err != nil {
			return err
		}

		if err := tx.Create(OAuthUser).Error; err != nil {
			return err
		}

		return nil
	})

	return OAuthUser, err
}
