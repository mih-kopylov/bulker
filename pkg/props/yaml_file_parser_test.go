package props

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestYamlFileParser_GetProperty(t *testing.T) {
	tests := []struct {
		name             string
		filename         string
		path             string
		expectedError    string
		expectedProperty *ParsedProperty
	}{
		{
			name: "string", filename: "testdata/data.yaml", path: "$.string", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      8,
				End:        19,
				Text:       "stringValue",
				Type:       PropertyTypeString,
				QuotesType: QuotesTypeNone,
			},
		},
		{
			name: "stringInQuotes", filename: "testdata/data.yaml", path: "$.stringInQuotes", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      36,
				End:        52,
				Text:       "stringInQuotes",
				Type:       PropertyTypeString,
				QuotesType: QuotesTypeDouble,
			},
		},
		{
			name: "number", filename: "testdata/data.yaml", path: "$.number", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      61,
				End:        63,
				Text:       "15",
				Type:       PropertyTypeNumber,
				QuotesType: QuotesTypeNone,
			},
		},
		{
			name: "numberNegative", filename: "testdata/data.yaml", path: "$.numberNegative", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      80,
				End:        83,
				Text:       "-15",
				Type:       PropertyTypeNumber,
				QuotesType: QuotesTypeNone,
			},
		},
		{
			name: "booleanTrue", filename: "testdata/data.yaml", path: "$.booleanTrue", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      97,
				End:        101,
				Text:       "true",
				Type:       PropertyTypeBool,
				QuotesType: QuotesTypeNone,
			},
		},
		{
			name: "booleanFalse", filename: "testdata/data.yaml", path: "$.booleanFalse", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      116,
				End:        121,
				Text:       "false",
				Type:       PropertyTypeBool,
				QuotesType: QuotesTypeNone,
			},
		},
		{
			name: "null", filename: "testdata/data.yaml", path: "$.null", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      128,
				End:        132,
				Text:       "null",
				Type:       PropertyTypeNull,
				QuotesType: QuotesTypeNone,
			},
		},
		{
			name: "tilda", filename: "testdata/data.yaml", path: "$.tilda", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      140,
				End:        141,
				Text:       "~",
				Type:       PropertyTypeNull,
				QuotesType: QuotesTypeNone,
			},
		},
		{
			name: "nested", filename: "testdata/data.yaml", path: "$.nested.string", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      160,
				End:        177,
				Text:       "nestedStringValue",
				Type:       PropertyTypeString,
				QuotesType: QuotesTypeNone,
			},
		},
		{
			name: "arrayString", filename: "testdata/data.yaml", path: "$.arrayString",
			expectedError: "middle path node kind is expected to be 8: kind=2, pos=12:3, path=$.arrayString",
		},
		{
			name: "array first element", filename: "testdata/data.yaml", path: "$.arrayString.0",
			expectedError: "middle path node kind is expected to be 4: kind=2, pos=12:3, path=$.arrayString",
		},
		{
			name: "array element by value", filename: "testdata/data.yaml", path: "$.arrayString.one",
			expectedError: "middle path node kind is expected to be 4: kind=2, pos=12:3, path=$.arrayString",
		},
		{
			name: "arrayObject", filename: "testdata/data.yaml", path: "$.arrayObject",
			expectedError: "middle path node kind is expected to be 8: kind=2, pos=15:3, path=$.arrayObject",
		},
		{
			name: "arrayObjectKey", filename: "testdata/data.yaml", path: "$.arrayObject.0.key",
			expectedError: "middle path node kind is expected to be 4: kind=2, pos=15:3, path=$.arrayObject",
		},
	}
	parser := YamlFileParser{}
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
