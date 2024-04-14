package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pauloo27/shurl/internal/server/api"
	"github.com/stretchr/testify/assert"
)

func TestErr(t *testing.T) {
	rr := httptest.NewRecorder()

	api.Err(rr, api.NotFoundErr, "Resource not found")
	assert.Equal(t, http.StatusNotFound, rr.Code)

	var responseBody map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&responseBody)
	assert.NoError(t, err)

	assert.Equal(t, string(api.NotFoundErr.Name), responseBody["error"])
	errorDetailMessage := responseBody["detail"].(map[string]interface{})["message"]
	assert.Equal(t, "Resource not found", errorDetailMessage)
}

type TestDetail struct {
	Name string `json:"name"`
}

func TestCreated(t *testing.T) {
	rr := httptest.NewRecorder()

	hello := TestDetail{
		Name: "Hello",
	}

	api.Created(rr, hello)
	assert.Equal(t, http.StatusCreated, rr.Code)

	var responseBody map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&responseBody)
	assert.NoError(t, err)

	assert.Equal(t, "Hello", responseBody["name"])
}
