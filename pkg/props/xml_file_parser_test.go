package props

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestXmlFileParser_GetProperty(t *testing.T) {
	tests := []struct {
		name             string
		filename         string
		path             string
		expectedError    string
		expectedProperty *ParsedProperty
	}{
		{
			name: "string", filename: "testdata/data.xml", path: "$.string", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      60,
				End:        71,
				Text:       "stringValue",
				Type:       PropertyTypeString,
				QuotesType: QuotesTypeNone,
			},
		},
		{
			name: "number", filename: "testdata/data.xml", path: "$.number", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      93,
				End:        95,
				Text:       "15",
				Type:       PropertyTypeString,
				QuotesType: QuotesTypeNone,
			},
		},
		{
			name: "numberNegative", filename: "testdata/data.xml", path: "$.numberNegative", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      125,
				End:        128,
				Text:       "-15",
				Type:       PropertyTypeString,
				QuotesType: QuotesTypeNone,
			},
		},
		{
			name: "booleanTrue", filename: "testdata/data.xml", path: "$.booleanTrue", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      163,
				End:        167,
				Text:       "true",
				Type:       PropertyTypeString,
				QuotesType: QuotesTypeNone,
			},
		},
		{
			name: "booleanFalse", filename: "testdata/data.xml", path: "$.booleanFalse", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      200,
				End:        205,
				Text:       "false",
				Type:       PropertyTypeString,
				QuotesType: QuotesTypeNone,
			},
		},
		{
			name: "null", filename: "testdata/data.xml", path: "$.null", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      231,
				End:        235,
				Text:       "null",
				Type:       PropertyTypeString,
				QuotesType: QuotesTypeNone,
			},
		},
		{
			name: "nested", filename: "testdata/data.xml", path: "$.nested.string", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      272,
				End:        289,
				Text:       "nestedStringValue",
				Type:       PropertyTypeString,
				QuotesType: QuotesTypeNone,
			},
		},
		{
			name: "array", filename: "testdata/data.xml", path: "$.duplicate", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      328,
				End:        338,
				Text:       "duplicate1",
				Type:       PropertyTypeString,
				QuotesType: QuotesTypeNone,
			},
		},
		{
			name: "empty", filename: "testdata/data.xml", path: "$.empty", expectedError: "",
			expectedProperty: &ParsedProperty{
				Start:      401,
				End:        401,
				Text:       "",
				Type:       PropertyTypeString,
				QuotesType: QuotesTypeNone,
			},
		},
	}
	parser := XmlFileParser{}
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
