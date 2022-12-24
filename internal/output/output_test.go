package output

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValueToMap(t *testing.T) {
	type a struct {
		Message string
	}
	result := valueToMap(a{"Hi"})
	assert.Equal(
		t, map[string]any{
			"message": "Hi",
		}, result,
	)
}

func TestValueToMap_CamelCase(t *testing.T) {
	type a struct {
		MessageCamelCase string
	}
	result := valueToMap(a{"Hi"})
	assert.Equal(
		t, map[string]any{
			"messageCamelCase": "Hi",
		}, result,
	)
}

func TestValueToMap_PrivateProperty(t *testing.T) {
	type a struct {
		messageCamelCase string
	}
	defer func() {
		rec := recover()
		assert.NotNilf(t, rec, "expected to panic")
	}()
	_ = valueToMap(a{"Hi"})
}
