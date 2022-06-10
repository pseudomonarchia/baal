package controller

import (
	"baal/lib/errorcode"
	"baal/lib/header"
	"baal/model"
	"baal/service"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// OAuth ...
type OAuth struct {
	Service *service.Services
}

// LoginURL get sso login url
func (o *OAuth) LoginURL(c *gin.Context) {
	query := &model.GoogleOAuthRequest{}
	_ = c.ShouldBindQuery(query)
	if _, ok := validator.New().Struct(query).(validator.ValidationErrors); ok {
		redirectURL, _ := url.Parse(c.GetHeader(header.Origin))
		ThrowErrorRedirect(c, redirectURL, errorcode.OAuthRequestQueryInvalid)
		return
	}

	OAuthState := o.Service.OAuth.NewState()
	expires := time.Now().Add(10 * time.Minute)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "oauthstate",
		Value:   OAuthState,
		Expires: expires,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "oauthredirect",
		Value:   query.Redirect,
		Expires: expires,
	})

	url := o.Service.OAuth.GetLoginURL(c.Request.Host, OAuthState)
	c.Redirect(http.StatusSeeOther, url)
}

// LoginCallBack get sso login callback
func (o *OAuth) LoginCallBack(c *gin.Context) {
	query := &model.GoogleOAuthResponse{}
	_ = c.ShouldBindQuery(query)
	redirectClient, errWithRedirect := c.Cookie("oauthredirect")
	OAuthState, errWithState := c.Cookie("oauthstate")
	redirectURL, err := url.Parse(redirectClient)
	rawQuery := url.Values{}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "oauthstate",
		Value:   "",
		Expires: time.Now(),
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "oauthredirect",
		Value:   "",
		Expires: time.Now(),
	})

	if errWithRedirect != nil {
		ThrowErrorRedirect(c, redirectURL, errorcode.ServerBasic)
		return
	}

	if _, ok := validator.New().Struct(query).(validator.ValidationErrors); ok {
		ThrowErrorRedirect(c, redirectURL, errorcode.OAuthResponseQueryInvalid)
		return
	}

	if errWithState != nil || query.State != OAuthState {
		ThrowErrorRedirect(c, redirectURL, errorcode.OAuthStateInvalid)
		return
	}

	token, err := o.Service.OAuth.GetToken(c.Request.Host, query.Code)
	if err != nil {
		ThrowErrorRedirect(c, redirectURL, errorcode.OAuthTokenInvalid)
		return
	}

	OAuthInfo, err := o.Service.OAuth.GetInfo(c.Request.Host, token)
	if err != nil {
		ThrowErrorRedirect(c, redirectURL, errorcode.OAuthTokenReject)
		return
	}

	userRef := &model.UserSchema{Email: OAuthInfo.Email}
	err = o.Service.User.GetByQuery(userRef)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ThrowErrorRedirect(c, redirectURL, errorcode.NotFoundOAuthUser)
		return
	}

	if !userRef.Enable {
		ThrowErrorRedirect(c, redirectURL, errorcode.UserDisabled)
		return
	}

	UID := o.Service.OAuth.GenerateUID()
	buf, _ := json.Marshal(token)
	OAuthRef := &model.OAuthTokenSchema{
		UID:       UID,
		UserID:    userRef.ID,
		TokenInfo: datatypes.JSON(buf),
		Provider:  model.OAuthProviderGoogle,
	}

	_ = o.Service.OAuth.SaveToken(OAuthRef)
	rawQuery.Add("code", UID)
	redirectURL.RawQuery = rawQuery.Encode()
	c.Redirect(http.StatusSeeOther, redirectURL.String())
}

// Token ...
func (o *OAuth) Token(c *gin.Context) {
	body := &model.TokenRequest{}
	_ = c.ShouldBindJSON(body)
	if _, ok := validator.New().Struct(body).(validator.ValidationErrors); ok {
		ThrowError(c, errorcode.LoginRequestBodyInvalid)
		return
	}

	var IP = c.ClientIP()
	var now = time.Now()
	var secret []byte
	var token *oauth2.Token
	var payload *model.JWTClaims
	var OAuthRef *model.OAuthTokenSchema
	var refreshRef *model.OAuthRefreshSchema

	switch body.GrantType {
	case model.GrantTypeFromCode:
		OAuthRef = &model.OAuthTokenSchema{UID: body.Code, Use: false}
		err := o.Service.OAuth.FindToken(OAuthRef, "Use")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ThrowError(c, errorcode.OAuthLoginFailed)
			return
		}

		token, _ = OAuthRef.UnmarshalToken()
		break
	case model.GrantTypeFromRefreshToken:
		auth := c.GetHeader(header.Authorization)
		payload = &model.JWTClaims{}

		if !payload.CheckBearerToken(auth) {
			ThrowError(c, errorcode.OAuthTokenFormatInvalid)
			return
		}

		if err := payload.DecodeToken(auth); err != nil {
			ThrowError(c, errorcode.OAuthTokenFormatInvalid)
			return
		}

		refreshRef = &model.OAuthRefreshSchema{Token: body.Code, OAuthUID: payload.Id}
		err := o.Service.OAuth.FindRefreshToken(refreshRef)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ThrowError(c, errorcode.OAuthTokenNotFound)
			return
		}

		if refreshRef.ExpiresAt.Before(now) || refreshRef.ExpiresAt.Equal(now) {
			o.Service.OAuth.CleanTokenResource(refreshRef)
			ThrowError(c, errorcode.OAuthTokenReject)
			return
		}

		if refreshRef.IP != IP {
			o.Service.OAuth.CleanTokenResource(refreshRef)
			ThrowError(c, errorcode.OAuthIssuedIPFailed)
			return
		}

		if !refreshRef.OAuthToken.User.Enable {
			ThrowError(c, errorcode.UserDisabled)
			return
		}

		secret = o.Service.OAuth.HashJSONToSHA(&refreshRef.OAuthToken.TokenInfo)
		if err := payload.ValidateToken(auth, secret); err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
			ThrowError(c, errorcode.OAuthTokenInvalid)
			return
		}

		token, _ = refreshRef.OAuthToken.UnmarshalToken()
		token, err = o.Service.OAuth.RefreshToken(c.Request.Host, token)
		if err != nil {
			ThrowError(c, errorcode.OAuthTokenReject)
			return
		}

		buf, _ := json.Marshal(token)
		OAuthRef = &model.OAuthTokenSchema{
			UID:       refreshRef.OAuthUID,
			UserID:    refreshRef.OAuthToken.UserID,
			TokenInfo: datatypes.JSON(buf),
			Provider:  model.OAuthProviderGoogle,
		}
	}

	secret = o.Service.OAuth.HashJSONToSHA(&OAuthRef.TokenInfo)
	refreshRef, _ = o.Service.OAuth.GenerateRefreshToken(OAuthRef, IP)
	payload = &model.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    c.Request.Host,
			Subject:   c.GetHeader(header.Origin),
			Audience:  strconv.Itoa(int(OAuthRef.UserID)),
			Id:        OAuthRef.UID,
			ExpiresAt: token.Expiry.Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
		},
	}

	c.JSON(http.StatusCreated, &model.TokenSchema{
		AccessToken:  payload.SignToken(secret),
		Expiry:       token.Expiry,
		TokenType:    token.TokenType,
		RefreshToken: refreshRef.Token,
	})
}
