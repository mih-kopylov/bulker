package props

import (
	"fmt"
	"github.com/buger/jsonparser"
)

type JsonFileParser struct {
}

func (j *JsonFileParser) GetProperty(fileBytes []byte, pp *ParsedPath) (*ParsedProperty, error) {
	value, dataType, offset, err := jsonparser.Get(fileBytes, pp.Elements...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	switch dataType {
	case jsonparser.NotExist:
		return nil, fmt.Errorf("no key found by path: path=%v", pp.String())
	case jsonparser.Unknown:
		fallthrough
	case jsonparser.Array:
		fallthrough
	case jsonparser.Object:
		return nil, fmt.Errorf("node expected to be a primitive: path=%v type=%v", pp.String(), dataType)
	case jsonparser.Null:
		return &ParsedProperty{
			offset - len(value),
			offset,
			string(value),
			PropertyTypeNull,
			QuotesTypeNone,
		}, nil
	case jsonparser.String:
		return &ParsedProperty{
			offset - len(value) - 1,
			offset - 1,
			string(value),
			PropertyTypeString,
			QuotesTypeDouble,
		}, nil
	case jsonparser.Boolean:
		return &ParsedProperty{
			offset - len(value),
			offset,
			string(value),
			PropertyTypeBool,
			QuotesTypeNone,
		}, nil
	case jsonparser.Number:
		return &ParsedProperty{
			offset - len(value),
			offset,
			string(value),
			PropertyTypeNumber,
			QuotesTypeNone,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported node type: path=%v type=%v", pp.String(), dataType)
	}
}
