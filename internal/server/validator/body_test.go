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

func TestValidateMissingBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/test", nil)
	res := httptest.NewRecorder()

	body, ok := validator.MustGetBody[SampleBody](res, req)
	assert.False(t, ok)
	assert.Zero(t, body)
}

func TestValidateInvalidData(t *testing.T) {
	rawBody := []byte(`{"name": "John Doe", "age": 1337}`)

	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(rawBody))
	res := httptest.NewRecorder()

	body, ok := validator.MustGetBody[SampleBody](res, req)
	assert.False(t, ok)
	assert.Equal(t, "John Doe", body.Name)
}

func TestValidJSONBody(t *testing.T) {
	rawBody := []byte(`{"name": "John Doe", "age": 30}`)

	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(rawBody))
	res := httptest.NewRecorder()

	body, ok := validator.MustGetBody[SampleBody](res, req)
	assert.True(t, ok)
	assert.Equal(t, "John Doe", body.Name)
}

type MockReader struct{}

func (r *MockReader) Read(p []byte) (int, error) {
	return 0, io.ErrClosedPipe
}

func TestClosedPipeBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/test", &MockReader{})
	res := httptest.NewRecorder()

	body, ok := validator.MustGetBody[SampleBody](res, req)
	assert.False(t, ok)
	assert.Zero(t, body)
}
