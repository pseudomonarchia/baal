package errorcode

import "net/http"

type errBasicType string
type errAuthorizeType errBasicType
type errForbiddenType errBasicType
type errParameterType errBasicType
type errNotFoundType errBasicType
type errRuntimeType errBasicType
type errServerType errBasicType

// Error code for client
const (
	Basic errBasicType = "000-00000"

	AuthorizeBasic          errAuthorizeType = "001-00000"
	OAuthStateInvalid       errAuthorizeType = "001-00001"
	OAuthTokenInvalid       errAuthorizeType = "001-00002"
	OAuthTokenReject        errAuthorizeType = "001-00003"
	OAuthLoginFailed        errAuthorizeType = "001-00004"
	OAuthTokenFormatInvalid errAuthorizeType = "001-00005"
	OAuthTokenNotFound      errAuthorizeType = "001-00006"
	OAuthIssuedIPFailed     errAuthorizeType = "001-00007"

	ForbiddenBasic errForbiddenType = "002-00000"
	UserDisabled   errForbiddenType = "002-00001"

	ParameterBasic            errParameterType = "003-00000"
	OAuthRequestQueryInvalid  errParameterType = "003-00001"
	OAuthResponseQueryInvalid errParameterType = "003-00002"
	LoginRequestBodyInvalid   errParameterType = "003-00003"

	NotFoundBasic     errNotFoundType = "004-30000"
	NotFoundOAuthUser errNotFoundType = "004-30001"

	RuntimeBasic errRuntimeType = "005-10000"

	ServerBasic errServerType = "006-10000"
)

// GetHTTPCode reply http code according to errorcode
func GetHTTPCode(code interface{}) int {
	switch code.(type) {
	case errBasicType:
		return http.StatusOK
	case errAuthorizeType:
		return http.StatusUnauthorized
	case errForbiddenType:
		return http.StatusForbidden
	case errParameterType:
		return http.StatusBadRequest
	case errNotFoundType:
		return http.StatusNotFound
	case errRuntimeType:
		return http.StatusOK
	case errServerType:
		return http.StatusInternalServerError
	default:
		return http.StatusOK
	}
}
