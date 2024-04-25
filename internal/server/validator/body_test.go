package validator_test

import (
	"bytes"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/pauloo27/shurl/internal/server/validator"
	"github.com/stretchr/testify/assert"
)

type SampleBody struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
	Age  int    `json:"age" validate:"required,min=1,max=150"`
}

func assertErrorMessage(t *testing.T, err *validator.APIBodyValidationError, message string) {
	if m, ok := err.Details.(map[string]string); ok {
		assert.Equal(t, message, m["message"])
	} else {
		t.Errorf("Error details is not a map[string]string")
	}
}

func assertFieldInvalid(t *testing.T, err *validator.APIBodyValidationError, field string, errorMsg string) {
	if errs, ok := err.Details.([]*validator.ValidationError); ok {
		found := false
		for _, e := range errs {
			if e.Field == field && e.Error == errorMsg {
				found = true
				break
			}
		}
		assert.True(t, found)
	} else {
		t.Errorf("Error details is not a []*validator.ValidationError")
	}
}

func TestValidateMissingBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/test", nil)

	body, err := validator.MustGetBody[SampleBody](req)

	assert.NotNil(t, err)
	assert.Zero(t, body)
	assert.Equal(t, "BAD_REQUEST", err.Error.Name)
	assert.Equal(t, 400, err.Error.StatusCode)
	assertErrorMessage(t, err, "EOF")
}

func TestValidateInvalidData(t *testing.T) {
	rawBody := []byte(`{"name": "John Doe", "age": 1337}`)

	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(rawBody))

	body, err := validator.MustGetBody[SampleBody](req)

	assert.NotNil(t, err)
	assert.Equal(t, SampleBody{Name: "John Doe", Age: 1337}, body)
	assertFieldInvalid(t, err, "age", "max 150")
}

func TestValidJSONBody(t *testing.T) {
	rawBody := []byte(`{"name": "John Doe", "age": 30}`)

	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(rawBody))

	body, err := validator.MustGetBody[SampleBody](req)
	assert.Nil(t, err)
	assert.Equal(t, SampleBody{Name: "John Doe", Age: 30}, body)
}

type MockReader struct{}

func (r *MockReader) Read(p []byte) (int, error) {
	return 0, io.ErrClosedPipe
}

func TestClosedPipeBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/test", &MockReader{})

	body, err := validator.MustGetBody[SampleBody](req)
	assert.NotNil(t, err)
	assert.Zero(t, body)
	assertErrorMessage(t, err, "io: read/write on closed pipe")
}
