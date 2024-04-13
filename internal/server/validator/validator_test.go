package validator_test

import (
	"testing"

	"github.com/pauloo27/shurl/internal/server/validator"
	"github.com/stretchr/testify/assert"
)

type SampleStruct struct {
	Name string `validate:"omitempty,min=3,max=20"`
	URL  string `validate:"required,http_url"`
	TTL  *int   `validate:"required"`
}

type SampleStructWithJSONTags struct {
	Name     string `json:"name,omitempty" validate:"omitempty,min=3,max=20"`
	URL      string `json:"url" validate:"required,http_url"`
	TTL      *int   `json:"ttl" validate:"required"`
	IgnoreMe string `json:"-"`
}

func TestValidData(t *testing.T) {
	ttl := 2
	data := SampleStruct{
		Name: "slug",
		URL:  "http://google.com",
		TTL:  &ttl,
	}

	errs := validator.Validate(data)
	assert.Empty(t, errs)
}

func TestValidDataWithMissingOptionals(t *testing.T) {
	ttl := 2
	data := SampleStruct{
		URL: "http://google.com",
		TTL: &ttl,
	}

	errs := validator.Validate(data)
	assert.Empty(t, errs)
}

func TestInvalidDataWithMissingRequireds(t *testing.T) {
	data := SampleStruct{
		Name: "john",
	}

	errs := validator.Validate(data)
	assert.NotEmpty(t, errs)
	assert.Len(t, errs, 2)

	errURL := errs[0]
	errTTL := errs[1]

	assert.Equal(t, "URL", errURL.Field)
	assert.Equal(t, "required", errURL.Error)

	assert.Equal(t, "TTL", errTTL.Field)
	assert.Equal(t, "required", errTTL.Error)
}

func TestInvalidDataWithErrParams(t *testing.T) {
	ttl := 2
	data := SampleStruct{
		Name: "g",
		URL:  "http://google.com",
		TTL:  &ttl,
	}

	errs := validator.Validate(data)
	assert.NotEmpty(t, errs)
	assert.Len(t, errs, 1)

	nameErr := errs[0]

	assert.Equal(t, "Name", nameErr.Field)
	assert.Equal(t, "min 3", nameErr.Error)
}

func TestInvalidDataWithJSONTags(t *testing.T) {
	data := SampleStructWithJSONTags{
		Name: "slug",
	}

	errs := validator.Validate(data)
	assert.NotEmpty(t, errs)
	assert.Len(t, errs, 2)

	errURL := errs[0]
	errTTL := errs[1]

	assert.Equal(t, "url", errURL.Field)
	assert.Equal(t, "required", errURL.Error)

	assert.Equal(t, "ttl", errTTL.Field)
	assert.Equal(t, "required", errTTL.Error)
}
