package controller_test

import (
	"baal/lib/errorcode"
	"baal/model"
	"baal/service/mocks"
	"baal/test"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func TestLoginURL(t *testing.T) {
	assert := assert.New(t)
	OAuthService := &mocks.OAuthFace{}
	r := test.MockSrvRoute(OAuthService)
	apiURL := &url.URL{Path: "/api/v1/oauth"}
	redirectURL := &url.URL{Scheme: "http", Host: "127.0.0.1"}

	state := "state_string"
	OAuthService.On("NewState").Return(state)
	OAuthService.On("GetLoginURL", "", state).Return("")

	t.Run("Does not carry query redirect return failure", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", apiURL.String(), nil)

		r.ServeHTTP(w, req)

		errorQuery := &url.Values{}
		errorQuery.Add("error_code", string(errorcode.OAuthRequestQueryInvalid))
		redirectURL.RawQuery = errorQuery.Encode()
		assert.Equal(http.StatusSeeOther, w.Code)
	})

	t.Run("Carry query redirect return success", func(t *testing.T) {
		query := &url.Values{}
		query.Add("redirect", redirectURL.String())
		apiURL.RawQuery = query.Encode()

		req, _ := http.NewRequest("GET", apiURL.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(http.StatusSeeOther, w.Code)
	})
}

func TestLoginCallBack(t *testing.T) {
	assert := assert.New(t)
	OAuthService := &mocks.OAuthFace{}
	userService := &mocks.UserFace{}
	state := "state_string"
	OAuthService.On("NewState").Return(state)
	OAuthService.On("GetLoginURL", "", state).Return("")

	r := test.MockSrvRoute(OAuthService, userService)
	apiURL := &url.URL{Path: "/api/v1/oauth/callback"}
	redirectURL := url.URL{Scheme: "http", Host: "127.0.0.1"}
	preAPIURL := &url.URL{
		Path:     "/api/v1/oauth",
		RawQuery: fmt.Sprintf("redirect=%s", redirectURL.String()),
	}

	t.Run("Not get login url return failure", func(t *testing.T) {
		req, _ := http.NewRequest("GET", apiURL.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(http.StatusSeeOther, w.Code)
	})

	t.Run("OAuth query invalid return failure", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", preAPIURL.String(), nil)
		r.ServeHTTP(w, req)
		assert.Equal(http.StatusSeeOther, w.Code)

		req, _ = http.NewRequest("GET", apiURL.String(), nil)
		for _, cookie := range w.Result().Cookies() {
			req.AddCookie(cookie)
		}

		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		errorQuery := &url.Values{}
		errorQuery.Add("error_code", string(errorcode.OAuthResponseQueryInvalid))
		redirectURL.RawQuery = errorQuery.Encode()
		actualURL, _ := w.Result().Location()
		assert.Equal(http.StatusSeeOther, w.Code)
		assert.Equal(redirectURL.String(), actualURL.String())
	})

	t.Run("OAuth state invalid return failure", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", preAPIURL.String(), nil)
		r.ServeHTTP(w, req)
		assert.Equal(http.StatusSeeOther, w.Code)

		OAuthQuery := &url.Values{}
		OAuthQuery.Add("state", "state")
		OAuthQuery.Add("code", "code")
		OAuthQuery.Add("scope", "scope")
		OAuthQuery.Add("authuser", "authuser")
		OAuthQuery.Add("prompt", "prompt")
		apiURL.RawQuery = OAuthQuery.Encode()

		req, _ = http.NewRequest("GET", apiURL.String(), nil)
		for _, cookie := range w.Result().Cookies() {
			req.AddCookie(cookie)
		}

		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		errorQuery := &url.Values{}
		errorQuery.Add("error_code", string(errorcode.OAuthStateInvalid))
		redirectURL.RawQuery = errorQuery.Encode()
		actualURL, _ := w.Result().Location()
		assert.Equal(http.StatusSeeOther, w.Code)
		assert.Equal(redirectURL.String(), actualURL.String())
	})

	t.Run("OAuth get token error return failure", func(t *testing.T) {
		codeStr := "code"
		OAuthService.On("GetToken", "", codeStr).Return(nil, errors.New("Token Error")).Once()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", preAPIURL.String(), nil)
		r.ServeHTTP(w, req)
		assert.Equal(http.StatusSeeOther, w.Code)

		OAuthQuery := &url.Values{}
		OAuthQuery.Add("state", state)
		OAuthQuery.Add("code", codeStr)
		OAuthQuery.Add("scope", "scope")
		OAuthQuery.Add("authuser", "authuser")
		OAuthQuery.Add("prompt", "prompt")
		apiURL.RawQuery = OAuthQuery.Encode()

		req, _ = http.NewRequest("GET", apiURL.String(), nil)
		for _, cookie := range w.Result().Cookies() {
			req.AddCookie(cookie)
		}

		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		errorQuery := &url.Values{}
		errorQuery.Add("error_code", string(errorcode.OAuthTokenInvalid))
		redirectURL.RawQuery = errorQuery.Encode()
		actualURL, _ := w.Result().Location()
		assert.Equal(http.StatusSeeOther, w.Code)
		assert.Equal(redirectURL.String(), actualURL.String())
	})

	t.Run("Get OAuth info error return failure", func(t *testing.T) {
		codeStr := "code"
		OAuthToken := &oauth2.Token{}
		OAuthService.On("GetToken", "", codeStr).Return(OAuthToken, nil).Once()
		OAuthService.On("GetInfo", "", OAuthToken).Return(nil, errors.New("Token Info Error")).Once()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", preAPIURL.String(), nil)
		r.ServeHTTP(w, req)
		assert.Equal(http.StatusSeeOther, w.Code)

		OAuthQuery := &url.Values{}
		OAuthQuery.Add("state", state)
		OAuthQuery.Add("code", codeStr)
		OAuthQuery.Add("scope", "scope")
		OAuthQuery.Add("authuser", "authuser")
		OAuthQuery.Add("prompt", "prompt")
		apiURL.RawQuery = OAuthQuery.Encode()

		req, _ = http.NewRequest("GET", apiURL.String(), nil)
		for _, cookie := range w.Result().Cookies() {
			req.AddCookie(cookie)
		}

		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		errorQuery := &url.Values{}
		errorQuery.Add("error_code", string(errorcode.OAuthTokenReject))
		redirectURL.RawQuery = errorQuery.Encode()
		actualURL, _ := w.Result().Location()
		assert.Equal(http.StatusSeeOther, w.Code)
		assert.Equal(redirectURL.String(), actualURL.String())
	})

	t.Run("Not found user return failure", func(t *testing.T) {
		codeStr := "code"
		OAuthToken := &oauth2.Token{}
		userData := &model.UserSchema{}
		OAuthService.On("GetToken", "", codeStr).Return(OAuthToken, nil).Once()
		OAuthService.On("GetInfo", "", OAuthToken).Return(&model.GoogleOAuthUserInfo{}, nil).Once()
		userService.On("GetByQuery", userData).Return(gorm.ErrRecordNotFound).Once()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", preAPIURL.String(), nil)
		r.ServeHTTP(w, req)
		assert.Equal(http.StatusSeeOther, w.Code)

		OAuthQuery := &url.Values{}
		OAuthQuery.Add("state", state)
		OAuthQuery.Add("code", codeStr)
		OAuthQuery.Add("scope", "scope")
		OAuthQuery.Add("authuser", "authuser")
		OAuthQuery.Add("prompt", "prompt")
		apiURL.RawQuery = OAuthQuery.Encode()

		req, _ = http.NewRequest("GET", apiURL.String(), nil)
		for _, cookie := range w.Result().Cookies() {
			req.AddCookie(cookie)
		}

		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		errorQuery := &url.Values{}
		errorQuery.Add("error_code", string(errorcode.NotFoundOAuthUser))
		redirectURL.RawQuery = errorQuery.Encode()
		actualURL, _ := w.Result().Location()
		assert.Equal(http.StatusSeeOther, w.Code)
		assert.Equal(redirectURL.String(), actualURL.String())
	})

	t.Run("Other conditions represent success", func(t *testing.T) {
		codeStr := "code"
		UID := "UID"
		OAuthToken := &oauth2.Token{}
		userData := &model.UserSchema{}
		buf, _ := json.Marshal(OAuthToken)
		OAuthData := &model.OAuthTokenSchema{
			UID:       UID,
			UserID:    userData.ID,
			TokenInfo: datatypes.JSON(buf),
			Provider:  model.OAuthProviderGoogle,
		}

		OAuthService.On("GetToken", "", codeStr).Return(OAuthToken, nil).Once()
		OAuthService.On("GetInfo", "", OAuthToken).Return(&model.GoogleOAuthUserInfo{}, nil).Once()
		userService.On("GetByQuery", userData).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(0).(*model.UserSchema)
			arg.Enable = true
		})

		OAuthService.On("GenerateUID").Return(UID).Once()
		OAuthService.On("SaveToken", OAuthData).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(0).(*model.OAuthTokenSchema)
			arg.Use = true
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", preAPIURL.String(), nil)
		r.ServeHTTP(w, req)
		assert.Equal(http.StatusSeeOther, w.Code)

		OAuthQuery := &url.Values{}
		OAuthQuery.Add("state", state)
		OAuthQuery.Add("code", codeStr)
		OAuthQuery.Add("scope", "scope")
		OAuthQuery.Add("authuser", "authuser")
		OAuthQuery.Add("prompt", "prompt")
		apiURL.RawQuery = OAuthQuery.Encode()

		req, _ = http.NewRequest("GET", apiURL.String(), nil)
		for _, cookie := range w.Result().Cookies() {
			req.AddCookie(cookie)
		}

		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		successQuery := &url.Values{}
		successQuery.Add("code", OAuthData.UID)
		redirectURL.RawQuery = successQuery.Encode()
		actualURL, _ := w.Result().Location()
		assert.Equal(http.StatusSeeOther, w.Code)
		assert.Equal(redirectURL.String(), actualURL.String())
	})
}
