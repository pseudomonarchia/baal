package controllers_test

import (
	"baal/test"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type response struct {
	Services string `json:"services"`
}

func TestIndexStatus(t *testing.T) {
	res, _ := json.Marshal(&response{"alive"})
	w := httptest.NewRecorder()
	r := test.MockRouter()
	req, _ := http.NewRequest("GET", "/api/v1/", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), string(res))
}
