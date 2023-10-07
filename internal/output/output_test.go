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

func TestValueKeys_KeysInSourceOrder(t *testing.T) {
	type a struct {
		J string
		W string
		A string
		Q string
		F string
		Z string
		O string
		C string
		M string
		H string
		Y string
		V string
		T string
		S string
		D string
		B string
		K string
		I string
		P string
		E string
		X string
		N string
		G string
		L string
		U string
		R string
	}
	keys := valueKeys(a{})

	assert.Equal(
		t, []string{
			"j", "w", "a", "q", "f", "z", "o", "c", "m", "h", "y", "v", "t", "s", "d", "b", "k",
			"i", "p", "e", "x", "n", "g", "l", "u", "r",
		},
		keys,
	)
}
