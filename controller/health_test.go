package controller_test

import (
	"baal/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	expected := "\"ok!\""
	w := httptest.NewRecorder()
	r := test.MockSrvRoute()

	req, _ := http.NewRequest("GET", "/health", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expected, w.Body.String())
}
