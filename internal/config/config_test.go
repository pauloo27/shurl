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

func TestInvalidConfigPath(t *testing.T) {
	cfg, err := config.LoadConfigFromFile("invalid")
	assert.Nil(t, cfg)
	assert.Error(t, err)
}

func TestLoadInvalidYamlData(t *testing.T) {
	cfg, err := config.LoadConfigFromData([]byte("x=x+1"))
	assert.Nil(t, cfg)
	assert.Error(t, err)
}

func TestLoadConfig(t *testing.T) {
	cfg, err := config.LoadConfigFromData([]byte("x=x+1"))
	assert.Nil(t, cfg)
	assert.Error(t, err)
}

func TestLoadEmptyConfig(t *testing.T) {
	cfg, err := config.LoadConfigFromData([]byte(""))
	assert.NotNil(t, cfg)
	assert.NoError(t, err)
	mustHaveNoNilPointers(t, *cfg)
}

func TestLoadConfigWithPublicAPIKey(t *testing.T) {
	// public with api key? nem a pau, juvenal
	cfg, err := config.LoadConfigFromData([]byte("public: { apiKey: true }"))
	assert.Nil(t, cfg)
	assert.Error(t, err)
}

func TestLoadDefaultConfig(t *testing.T) {
	cfg, err := config.LoadConfigFromFile(defaultConfigPath)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	mustHaveNoZeroValue(t, cfg)
}

func mustHaveNoNilPointers(t *testing.T, cfg config.Config) {
	valType := reflect.TypeOf(cfg)
	valVal := reflect.ValueOf(cfg)

	for i := 0; i < valType.NumField(); i++ {
		field := valType.Field(i)

		// skip fields ignored by yaml
		yamlTag := field.Tag.Get("yaml")
		if yamlTag == "-" {
			continue
		}

		fieldVal := valVal.Field(i).Interface()
		assert.NotNil(t, fieldVal, "Field %s must not be nil", field.Name)
	}
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
