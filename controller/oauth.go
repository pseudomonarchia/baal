package controller

import (
	"baal/lib/errorcode"
	"baal/model"
	"baal/service"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
		ec := errorcode.OAuthRequestQueryInvalid
		code := errorcode.GetHTTPCode(ec)
		c.JSON(code, model.ErrorResponse{ErrorCode: ec})
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
	c.JSON(http.StatusOK, gin.H{"url": url})
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
		ec := errorcode.ServerBasic
		code := errorcode.GetHTTPCode(ec)
		c.JSON(code, model.ErrorResponse{ErrorCode: ec})
		return
	}

	if _, ok := validator.New().Struct(query).(validator.ValidationErrors); ok {
		rawQuery.Add("error_code", string(errorcode.OAuthResponseQueryInvalid))
		redirectURL.RawQuery = rawQuery.Encode()
		c.Redirect(http.StatusSeeOther, redirectURL.String())
		return
	}

	if errWithState != nil || query.State != OAuthState {
		rawQuery.Add("error_code", string(errorcode.OAuthStateInvalid))
		redirectURL.RawQuery = rawQuery.Encode()
		c.Redirect(http.StatusSeeOther, redirectURL.String())
		return
	}

	token, err := o.Service.OAuth.GetToken(c.Request.Host, query.Code)
	if err != nil {
		rawQuery.Add("error_code", string(errorcode.OAuthTokenInvalid))
		redirectURL.RawQuery = rawQuery.Encode()
		c.Redirect(http.StatusSeeOther, redirectURL.String())
		return
	}

	OAuthInfo, err := o.Service.OAuth.GetInfo(c.Request.Host, token)
	if err != nil {
		rawQuery.Add("error_code", string(errorcode.OAuthTokenReject))
		redirectURL.RawQuery = rawQuery.Encode()
		c.Redirect(http.StatusSeeOther, redirectURL.String())
		return
	}

	user, notExist := o.Service.User.GetByQuery(&model.UserSchema{Email: OAuthInfo.Email})
	if notExist {
		rawQuery.Add("error_code", string(errorcode.NotFoundOAuthUser))
		redirectURL.RawQuery = rawQuery.Encode()
		c.Redirect(http.StatusSeeOther, redirectURL.String())
		return
	}

	OAuthRef, err := o.Service.OAuth.SaveToken(user.ID, token)
	rawQuery.Add("code", OAuthRef.UID)
	redirectURL.RawQuery = rawQuery.Encode()
	c.Redirect(http.StatusSeeOther, redirectURL.String())
}
