package props

import (
	"errors"
	"strings"
)

var (
	ErrPropertyNotFound        = errors.New("property not found")
	ErrMalformedPathExpressoin = errors.New("malformed path expression")
	ErrUnsupportedFileFormat   = errors.New("unsupported file format")
)

// ParsedPath describes steps of searching for a property in a file content
type ParsedPath struct {
	Elements []string
}

func (p *ParsedPath) String() string {
	return "$." + strings.Join(p.Elements, ".")
}

// ParsedProperty describes a found property
type ParsedProperty struct {
	// Start cursor position of the first character of the value. Starts from 0
	Start int
	// End cursor position of the next character after the last one of the value
	End int
	// Text text representation of the value
	Text string
	// Type describes type of the value
	Type PropertyType
	// QuotesType describes used type of quotes
	QuotesType QuotesType
}

// FileParser is a base interface that works with file content and returns a property based on provided path
type FileParser interface {
	GetProperty(fileBytes []byte, pp *ParsedPath) (*ParsedProperty, error)
}

type PropertyType int

const (
	PropertyTypeString PropertyType = 1 << iota
	PropertyTypeNumber
	PropertyTypeBool
	PropertyTypeNull
)

type QuotesType int

const (
	QuotesTypeNone QuotesType = 1 << iota
	QuotesTypeSingular
	QuotesTypeDouble
)
