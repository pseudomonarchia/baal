package controller

import (
	"baal/lib/errorcode"
	"baal/model"
	"baal/service"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// Controllers represents a global controllers struct
type Controllers struct {
	Index  *Index
	OAuth  *OAuth
	Health *Health
}

// New return all controller
func New(s *service.Services) *Controllers {
	return &Controllers{
		Index:  &Index{s},
		OAuth:  &OAuth{s},
		Health: &Health{s},
	}
}

// ThrowError ...
func ThrowError(c *gin.Context, ErrCode interface{}) {
	code := errorcode.GetHTTPCode(ErrCode)
	c.JSON(code, model.ErrorResponse{ErrorCode: ErrCode})
}

// ThrowErrorRedirect ...
func ThrowErrorRedirect(c *gin.Context, redirect *url.URL, ErrCode interface{}) {
	rawQuery := url.Values{}
	rawQuery.Add("error_code", fmt.Sprintf("%v", ErrCode))
	redirect.RawQuery = rawQuery.Encode()
	c.Redirect(http.StatusSeeOther, redirect.String())
}
