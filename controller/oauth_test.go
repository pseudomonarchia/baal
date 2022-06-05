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
	"golang.org/x/oauth2"
)

func TestLoginURL(t *testing.T) {
	assert := assert.New(t)
	OAuthService := &mocks.OAuthFace{}
	r := test.MockSrvRoute(OAuthService)
	apiURL := &url.URL{Path: "/api/v1/login"}
	redirectURL := &url.URL{Scheme: "http", Host: "127.0.0.1"}

	state := "state_string"
	OAuthService.On("NewState").Return(state)
	OAuthService.On("GetLoginURL", "", state).Return("")

	t.Run("Does not carry query redirect return failure", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", apiURL.String(), nil)
		e := errorcode.OAuthRequestQueryInvalid
		resStatus := errorcode.GetHTTPCode(e)
		expectedJSON, _ := json.Marshal(model.ErrorResponse{ErrorCode: e})

		r.ServeHTTP(w, req)
		assert.Equal(resStatus, w.Code)
		assert.Equal(string(expectedJSON), w.Body.String())
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
	apiURL := &url.URL{Path: "/api/v1/login/callback"}
	redirectURL := url.URL{Scheme: "http", Host: "127.0.0.1"}
	preAPIURL := &url.URL{
		Path:     "/api/v1/login",
		RawQuery: fmt.Sprintf("redirect=%s", redirectURL.String()),
	}

	t.Run("Not get login url return failure", func(t *testing.T) {
		req, _ := http.NewRequest("GET", apiURL.String(), nil)
		e := errorcode.ServerBasic
		expectedCode := errorcode.GetHTTPCode(e)
		expectedJSON, _ := json.Marshal(model.ErrorResponse{ErrorCode: e})

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(expectedCode, w.Code)
		assert.Equal(string(expectedJSON), w.Body.String())
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
		OAuthService.On("GetToken", "", codeStr).Return(OAuthToken, nil).Once()
		OAuthService.On("GetInfo", "", OAuthToken).Return(&model.GoogleOAuthUserInfo{}, nil).Once()
		userService.On("GetByQuery", &model.UserSchema{}).Return(nil, true).Once()
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
		OAuthToken := &oauth2.Token{}
		userData := &model.UserSchema{}
		OAuthData := &model.OAuthSchema{}
		OAuthService.On("GetToken", "", codeStr).Return(OAuthToken, nil).Once()
		OAuthService.On("GetInfo", "", OAuthToken).Return(&model.GoogleOAuthUserInfo{}, nil).Once()
		userService.On("GetByQuery", userData).Return(userData, false).Once()
		OAuthService.On("SaveToken", userData.ID, OAuthToken).Return(OAuthData, nil)
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
