package props

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestJsonFileParser_GetProperty(t *testing.T) {
	tests := []struct {
		name             string
		filename         string
		path             string
		expectedError    string
		expectedProperty *ParsedProperty
	}{
		{
			name: "string", filename: "testdata/data.json", path: "$.string", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      15,
				End:        26,
				Text:       "stringValue",
				Type:       PropertyTypeString,
				QuotesType: QuotesTypeDouble,
			},
		},
		{
			name: "number", filename: "testdata/data.json", path: "$.number", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      41,
				End:        43,
				Text:       "15",
				Type:       PropertyTypeNumber,
				QuotesType: QuotesTypeNone,
			},
		},
		{
			name: "numberNegative", filename: "testdata/data.json", path: "$.numberNegative", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      65,
				End:        68,
				Text:       "-15",
				Type:       PropertyTypeNumber,
				QuotesType: QuotesTypeNone,
			},
		},
		{
			name: "booleanTrue", filename: "testdata/data.json", path: "$.booleanTrue", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      87,
				End:        91,
				Text:       "true",
				Type:       PropertyTypeBool,
				QuotesType: QuotesTypeNone,
			},
		},
		{
			name: "booleanFalse", filename: "testdata/data.json", path: "$.booleanFalse", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      111,
				End:        116,
				Text:       "false",
				Type:       PropertyTypeBool,
				QuotesType: QuotesTypeNone,
			},
		},
		{
			name: "null", filename: "testdata/data.json", path: "$.null", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      128,
				End:        132,
				Text:       "null",
				Type:       PropertyTypeNull,
				QuotesType: QuotesTypeNone,
			},
		},
		{
			name: "nested", filename: "testdata/data.json", path: "$.nested.string", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      163,
				End:        180,
				Text:       "nestedStringValue",
				Type:       PropertyTypeString,
				QuotesType: QuotesTypeDouble,
			},
		},
		{
			name: "array", filename: "testdata/data.json", path: "$.arrayString",
			expectedError: "node expected to be a primitive: path=$.arrayString type=array",
		},
		{
			name: "array first element", filename: "testdata/data.json", path: "$.arrayString.0",
			expectedError: "failed to parse json: Key path not found",
		},
		{
			name: "array element by value", filename: "testdata/data.json", path: "$.arrayString.one",
			expectedError: "failed to parse json: Key path not found",
		},
		{
			name: "array", filename: "testdata/data.json", path: "$.arrayObject",
			expectedError: "node expected to be a primitive: path=$.arrayObject type=array",
		},
		{
			name: "array", filename: "testdata/data.json", path: "$.arrayObject.0.key",
			expectedError: "failed to parse json: Key path not found",
		},
	}
	parser := JsonFileParser{}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				bytes, err := os.ReadFile(test.filename)
				if !assert.NoError(t, err) {
					return
				}

				parsedPath, err := ParsePath(test.path)
				if !assert.NoError(t, err) {
					return
				}

				property, err := parser.GetProperty(bytes, parsedPath)
				if test.expectedError != "" {
					assert.Error(t, err)
					assert.Equal(t, test.expectedError, err.Error())
				} else {
					if !assert.NoError(t, err) {
						return
					}
					assert.Equal(t, test.expectedProperty, property)
				}
			},
		)
	}
}
