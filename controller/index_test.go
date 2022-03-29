package controller_test

import (
	"baal/test"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexStatus(t *testing.T) {
	expectedJSON, _ := json.Marshal(struct {
		Services string `json:"services"`
	}{
		"alive",
	})

	w := httptest.NewRecorder()
	r := test.MockSrvRoute()
	req, _ := http.NewRequest("GET", "/api/v1/", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(expectedJSON), w.Body.String())
}
