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
