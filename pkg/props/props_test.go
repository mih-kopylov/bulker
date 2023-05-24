package props

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestGetPropertyFromFile(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		path     string
		expected string
	}{
		{name: "json", fileName: "testdata/data.json", path: "$.nested.string", expected: "nestedStringValue"},
		{name: "xml", fileName: "testdata/data.xml", path: "$.nested.string", expected: "nestedStringValue"},
		{name: "yaml", fileName: "testdata/data.yaml", path: "$.nested.string", expected: "nestedStringValue"},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				prop, err := GetPropertyFromFile(test.fileName, test.path)
				if assert.NoError(t, err) {
					assert.Equal(t, test.expected, prop)
				}
			},
		)
	}
}

func TestGetFileParser(t *testing.T) {
	tests := []struct {
		name          string
		fileName      string
		expectedError bool
		expectedType  string
	}{
		{name: "no extension", fileName: "examplexml", expectedError: true},
		{name: "xml", fileName: "example.xml", expectedError: false, expectedType: "*props.XmlFileParser"},
		{name: "xmli", fileName: "example.xmli", expectedError: true},
		{
			name: "double extension", fileName: "example.json.xml", expectedError: false,
			expectedType: "*props.XmlFileParser",
		},
		{name: "just extension", fileName: ".xml", expectedError: false, expectedType: "*props.XmlFileParser"},
		{name: "extension without dot", fileName: "xml", expectedError: true},
		{name: "json", fileName: "example.json", expectedError: false, expectedType: "*props.JsonFileParser"},
		{name: "yml", fileName: "example.yml", expectedError: false, expectedType: "*props.YamlFileParser"},
		{name: "yaml", fileName: "example.yaml", expectedError: false, expectedType: "*props.YamlFileParser"},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				parser, err := getFileParser(test.fileName)
				if test.expectedError {
					assert.Error(t, err)
				} else {
					if assert.NoError(t, err) {
						typeOf := reflect.TypeOf(parser)
						assert.Equal(t, test.expectedType, typeOf.String())
					}
				}
			},
		)
	}
}

func TestParsePath(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		expectedError bool
		expected      *ParsedPath
	}{
		{
			name: "one level", path: "$.value", expectedError: false,
			expected: &ParsedPath{Elements: []string{"value"}},
		},
		{
			name: "two levels", path: "$.value.subvalue", expectedError: false,
			expected: &ParsedPath{Elements: []string{"value", "subvalue"}},
		},
		{
			name: "array", path: "$.value[0].subvalue", expectedError: false,
			expected: &ParsedPath{Elements: []string{"value[0]", "subvalue"}},
		},
		{
			name: "array", path: "$.value.[0].subvalue", expectedError: false,
			expected: &ParsedPath{Elements: []string{"value", "[0]", "subvalue"}},
		},
		{
			name: "no root", path: "value.subvalue", expectedError: true,
		},
		{
			name: "slashes", path: "/value/subvalue", expectedError: true,
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				actual, err := ParsePath(test.path)
				if test.expectedError {
					assert.Error(t, err)
				} else {
					if assert.NoError(t, err) {
						assert.Equal(t, test.expected, actual)
					}
				}
			},
		)
	}
}
