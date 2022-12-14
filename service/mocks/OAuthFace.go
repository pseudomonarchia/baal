// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	datatypes "gorm.io/datatypes"

	model "baal/model"

	oauth2 "golang.org/x/oauth2"
)

// OAuthFace is an autogenerated mock type for the OAuthFace type
type OAuthFace struct {
	mock.Mock
}

// CleanTokenResource provides a mock function with given fields: ref
func (_m *OAuthFace) CleanTokenResource(ref *model.OAuthRefreshSchema) error {
	ret := _m.Called(ref)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.OAuthRefreshSchema) error); ok {
		r0 = rf(ref)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindRefreshToken provides a mock function with given fields: ref, query
func (_m *OAuthFace) FindRefreshToken(ref *model.OAuthRefreshSchema, query ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, ref)
	_ca = append(_ca, query...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.OAuthRefreshSchema, ...interface{}) error); ok {
		r0 = rf(ref, query...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindToken provides a mock function with given fields: ref, query
func (_m *OAuthFace) FindToken(ref *model.OAuthTokenSchema, query ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, ref)
	_ca = append(_ca, query...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.OAuthTokenSchema, ...interface{}) error); ok {
		r0 = rf(ref, query...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GenerateHashFromUID provides a mock function with given fields: UID
func (_m *OAuthFace) GenerateHashFromUID(UID string) string {
	ret := _m.Called(UID)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(UID)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GenerateRefreshToken provides a mock function with given fields: ref, IP
func (_m *OAuthFace) GenerateRefreshToken(ref *model.OAuthTokenSchema, IP string) (*model.OAuthRefreshSchema, error) {
	ret := _m.Called(ref, IP)

	var r0 *model.OAuthRefreshSchema
	if rf, ok := ret.Get(0).(func(*model.OAuthTokenSchema, string) *model.OAuthRefreshSchema); ok {
		r0 = rf(ref, IP)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.OAuthRefreshSchema)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.OAuthTokenSchema, string) error); ok {
		r1 = rf(ref, IP)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateUID provides a mock function with given fields:
func (_m *OAuthFace) GenerateUID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetInfo provides a mock function with given fields: requestURL, token
func (_m *OAuthFace) GetInfo(requestURL string, token *oauth2.Token) (*model.GoogleOAuthUserInfo, error) {
	ret := _m.Called(requestURL, token)

	var r0 *model.GoogleOAuthUserInfo
	if rf, ok := ret.Get(0).(func(string, *oauth2.Token) *model.GoogleOAuthUserInfo); ok {
		r0 = rf(requestURL, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.GoogleOAuthUserInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *oauth2.Token) error); ok {
		r1 = rf(requestURL, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLoginURL provides a mock function with given fields: requestURL, state
func (_m *OAuthFace) GetLoginURL(requestURL string, state string) string {
	ret := _m.Called(requestURL, state)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(requestURL, state)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetToken provides a mock function with given fields: requestURL, code
func (_m *OAuthFace) GetToken(requestURL string, code string) (*oauth2.Token, error) {
	ret := _m.Called(requestURL, code)

	var r0 *oauth2.Token
	if rf, ok := ret.Get(0).(func(string, string) *oauth2.Token); ok {
		r0 = rf(requestURL, code)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth2.Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(requestURL, code)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HashJSONToSHA provides a mock function with given fields: r
func (_m *OAuthFace) HashJSONToSHA(r *datatypes.JSON) []byte {
	ret := _m.Called(r)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(*datatypes.JSON) []byte); ok {
		r0 = rf(r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

// NewState provides a mock function with given fields:
func (_m *OAuthFace) NewState() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// RefreshToken provides a mock function with given fields: requestURL, token
func (_m *OAuthFace) RefreshToken(requestURL string, token *oauth2.Token) (*oauth2.Token, error) {
	ret := _m.Called(requestURL, token)

	var r0 *oauth2.Token
	if rf, ok := ret.Get(0).(func(string, *oauth2.Token) *oauth2.Token); ok {
		r0 = rf(requestURL, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth2.Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *oauth2.Token) error); ok {
		r1 = rf(requestURL, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveToken provides a mock function with given fields: ref
func (_m *OAuthFace) SaveToken(ref *model.OAuthTokenSchema) error {
	ret := _m.Called(ref)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.OAuthTokenSchema) error); ok {
		r0 = rf(ref)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
