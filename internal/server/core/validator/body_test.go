package validator_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/pauloo27/shurl/internal/server/core/validator"
	"github.com/stretchr/testify/assert"
)

type SampleBody struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
	Age  int    `json:"age" validate:"required,min=1,max=150"`
}

func TestValidateMissingBody(t *testing.T) {
	ctx := newEchoCtx(nil)

	body, err := validator.MustBindAndValidate[SampleBody](ctx)
	assert.NotNil(t, err)
	assert.Zero(t, body)
	assert.Equal(t, "VALIDATION_ERROR", err.Error.Name)
	assert.Equal(t, 422, err.Error.StatusCode)
}

func TestValidateInvalidData(t *testing.T) {
	rawBody := `{"name": "John Doe", "age": 1337}`
	ctx := newEchoCtx(strings.NewReader(rawBody))
	body, err := validator.MustBindAndValidate[SampleBody](ctx)
	assert.NotNil(t, err)
	assert.Equal(t, SampleBody{Name: "John Doe", Age: 1337}, body)
	assertFieldInvalid(t, err, "age", "max 150")
}

func TestValidJSONBody(t *testing.T) {
	rawBody := `{"name": "John Doe", "age": 30}`
	ctx := newEchoCtx(strings.NewReader(rawBody))
	body, err := validator.MustBindAndValidate[SampleBody](ctx)
	assert.Nil(t, err)
	assert.Equal(t, SampleBody{Name: "John Doe", Age: 30}, body)
}

type MockReader struct{}

func (r *MockReader) Read(p []byte) (int, error) {
	return 0, io.ErrClosedPipe
}

func TestClosedPipeBody(t *testing.T) {
	ctx := newEchoCtx(&MockReader{})
	body, err := validator.MustBindAndValidate[SampleBody](ctx)
	assert.NotNil(t, err)
	assert.Zero(t, body)
	assertErrorMessage(t, err, "Invalid payload")
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

func newEchoCtx(body io.Reader) echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}
