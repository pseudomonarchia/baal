package controller

import (
	"baal/lib/errorcode"
	"baal/lib/header"
	"baal/model"
	"baal/service"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
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

	userRef, err := o.Service.User.GetByQuery(&model.UserSchema{Email: OAuthInfo.Email})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ThrowErrorRedirect(c, redirectURL, errorcode.NotFoundOAuthUser)
		return
	}

	if !userRef.Enable {
		ThrowErrorRedirect(c, redirectURL, errorcode.UserDisabled)
		return
	}

	UID := o.Service.OAuth.GenerateUID()
	OAuthRef, err := o.Service.OAuth.SaveToken(userRef.ID, UID, token)
	rawQuery.Add("code", OAuthRef.UID)
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

	switch body.GrantType {
	case model.GrantTypeFromCode:
		OAuthRef := o.Service.OAuth.FindToken(body.Code)
		if OAuthRef == nil {
			ThrowError(c, errorcode.OAuthLoginFailed)
			return
		}

		token, _ := OAuthRef.UnmarshalToken()
		payload := &jwt.StandardClaims{
			Issuer:    c.Request.Host,
			Subject:   c.GetHeader(header.Origin),
			Audience:  strconv.Itoa(int(OAuthRef.UserID)),
			Id:        OAuthRef.UID,
			ExpiresAt: token.Expiry.Unix(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
		}

		refreshRef, _ := o.Service.OAuth.GenerateRefreshToken(OAuthRef, c.ClientIP())
		res := &model.TokenSchema{
			AccessToken:  o.Service.OAuth.SignToken(payload, &OAuthRef.TokenInfo),
			Expiry:       token.Expiry,
			TokenType:    token.TokenType,
			RefreshToken: refreshRef.Token,
		}

		c.JSON(http.StatusCreated, res)
		break
	case model.GrantTypeFromRefreshToken:
		a := c.GetHeader(header.Authorization)
		t, ok := o.Service.OAuth.CheckTokenAndReplace(a)
		if !ok {
			ThrowError(c, errorcode.OAuthTokenFormatInvalid)
			return
		}

		payload, err := o.Service.OAuth.DecodeToken(t)
		if err != nil {
			ThrowError(c, errorcode.OAuthTokenFormatInvalid)
			return
		}

		refreshRef, err := o.Service.OAuth.FindRefreshTokenFormUID(body.Code, payload.Id)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ThrowError(c, errorcode.OAuthTokenNotFound)
			return
		}

		if IP := c.ClientIP(); refreshRef.IP != IP {
			ThrowError(c, errorcode.OAuthIssuedIPFailed)
			return
		}

		if !refreshRef.OAuthToken.User.Enable {
			ThrowError(c, errorcode.UserDisabled)
			return
		}

		validate, err := o.Service.OAuth.ValidateToken(t, &refreshRef.OAuthToken.TokenInfo)
		tokenExpired := errors.Is(err, jwt.ErrTokenExpired)
		if !validate && !tokenExpired {
			ThrowError(c, errorcode.OAuthTokenInvalid)
			return
		}

		token, _ := refreshRef.OAuthToken.UnmarshalToken()
		token, err = o.Service.OAuth.RefreshToken(c.Request.Host, token)
		if err != nil {
			ThrowError(c, errorcode.OAuthTokenReject)
			return
		}

		OAuthRef, _ := o.Service.OAuth.SaveToken(refreshRef.OAuthToken.UserID, refreshRef.OAuthUID, token)
		refreshRef, _ = o.Service.OAuth.GenerateRefreshToken(OAuthRef, c.ClientIP())

		payload.ExpiresAt = token.Expiry.Unix()
		payload.IssuedAt = time.Now().Unix()
		payload.NotBefore = time.Now().Unix()

		res := &model.TokenSchema{
			AccessToken:  o.Service.OAuth.SignToken(payload, &OAuthRef.TokenInfo),
			Expiry:       token.Expiry,
			TokenType:    token.TokenType,
			RefreshToken: refreshRef.Token,
		}

		c.JSON(http.StatusCreated, res)
	}
}
