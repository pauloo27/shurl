package config_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pauloo27/shurl/internal/config"
	"github.com/stretchr/testify/assert"
)

const (
	defaultConfigPath = "../../config.example.yaml"
)

var (
	mustBeUnset = map[string]bool{
		"Config.Public.APIKey": true,
	}
)

func TestLoadDefaultConfig(t *testing.T) {
	config, err := config.LoadConfig(defaultConfigPath)
	assert.Nil(t, err)
	assert.NotNil(t, config)

	mustHaveNoZeroValue(t, config)
}

func mustHaveNoZeroValue(t *testing.T, cfg *config.Config) {
	check(t, cfg, "Config")
}

func check(t *testing.T, val interface{}, path string) {
	valType := reflect.TypeOf(val)
	valVal := reflect.ValueOf(val)

	// de-reference pointers, will make our lives easier
	if valType.Kind() == reflect.Ptr {
		valVal = valVal.Elem()
		val = valVal.Interface()
		valType = reflect.TypeOf(val)
	}

	switch valType.Kind() {
	case reflect.String, reflect.Int, reflect.Bool:
		if mustBeUnset[path] {
			assert.Zerof(t, val, "Value for %s must BE zero, but was not", path)
		} else {
			assert.NotZerof(t, val, "Value for %s must NOT be zero, but was", path)
		}
	case reflect.Struct:
		for i := 0; i < valType.NumField(); i++ {
			field := valType.Field(i)

			// skip fields ignored by yaml
			yamlTag := field.Tag.Get("yaml")
			if yamlTag == "-" {
				continue
			}

			fieldPath := fmt.Sprintf("%s.%s", path, field.Name)
			check(t, valVal.Field(i).Interface(), fieldPath)
		}
	case reflect.Map:
		assert.Equal(t, reflect.String, valType.Key().Kind(), "Map key type must be string")
		for _, key := range valVal.MapKeys() {
			check(t, valVal.MapIndex(key).Interface(), fmt.Sprintf("%s[%s]", path, key.String()))
		}
	case reflect.Slice:
		for i := 0; i < valVal.Len(); i++ {
			check(t, valVal.Index(i).Interface(), fmt.Sprintf("%s[%d]", path, i))
		}
	default:
		panic(fmt.Sprintf("Unexpected field type %s", valType.Kind()))
	}
}
