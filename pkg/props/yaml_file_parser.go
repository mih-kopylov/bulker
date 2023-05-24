package props

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
)

type YamlFileParser struct {
}

func (p *YamlFileParser) GetProperty(fileBytes []byte, pp *ParsedPath) (*ParsedProperty, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(fileBytes))
	var documentNode yaml.Node
	err := decoder.Decode(&documentNode)
	if err != nil {
		return nil, fmt.Errorf("failed to decode yaml: %w", err)
	}

	if documentNode.Kind != yaml.DocumentNode {
		return nil, fmt.Errorf(
			"parsed documented expected to be a document node: kind=%v, pos=%v:%v",
			documentNode.Kind, documentNode.Line, documentNode.Column,
		)
	}

	node := documentNode.Content[0]
	if node.Kind != yaml.MappingNode {
		return nil, fmt.Errorf(
			"root yaml node is expected to be a map: kind=%v, pos=%v:%v", node.Kind, node.Line,
			node.Column,
		)
	}

	for elementIndex, element := range pp.Elements[:len(pp.Elements)-1] {
		//searching for a next key that should reference a map
		node, err = p.getKeyWithValueOfType(
			node, element, yaml.MappingNode,
			&ParsedPath{pp.Elements[0 : elementIndex+1]},
		)
		if err != nil {
			return nil, err
		}
	}
	//searching for the last path element that should reference a scalar
	node, err = p.getKeyWithValueOfType(
		node, pp.Elements[len(pp.Elements)-1], yaml.ScalarNode,
		&ParsedPath{pp.Elements},
	)
	if err != nil {
		return nil, err
	}

	quotesType := QuotesTypeNone
	if node.Style == yaml.DoubleQuotedStyle {
		quotesType = QuotesTypeDouble
	} else if node.Style == yaml.SingleQuotedStyle {
		quotesType = QuotesTypeSingular
	}

	start := getBytesLengthInLines(fileBytes, node.Line-1) + node.Column - 1
	end := start + len(node.Value)
	if quotesType != QuotesTypeNone {
		// adding length of quotes
		end += 2
	}

	var propertyType PropertyType
	switch node.Tag {
	case "!!str":
		propertyType = PropertyTypeString
	case "!!null":
		propertyType = PropertyTypeNull
	case "!!int":
		fallthrough
	case "!!float":
		propertyType = PropertyTypeNumber
	case "!!bool":
		propertyType = PropertyTypeBool
	default:
		return nil, errors.New("unsupported yaml tag: " + node.Tag)
	}

	return &ParsedProperty{
		start,
		end,
		node.Value,
		propertyType,
		quotesType,
	}, nil
}

func getBytesLengthInLines(fileBytes []byte, lines int) int {
	if lines == 0 {
		return 0
	}

	nIndex := bytes.IndexByte(fileBytes, '\n')
	rIndex := bytes.IndexByte(fileBytes, '\r')
	var lineSeparator []byte
	if rIndex < 0 {
		lineSeparator = []byte{'\n'}
	} else if nIndex < 0 {
		lineSeparator = []byte{'\r'}
	} else {
		lineSeparator = []byte{'\r', '\n'}
	}

	scanner := bufio.NewScanner(bytes.NewReader(fileBytes))
	scanner.Split(
		func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			if atEOF && len(data) == 0 {
				return 0, nil, nil
			}
			if i := bytes.Index(data, lineSeparator); i >= 0 {
				// We have a full newline-terminated line.
				return i + 1, data[0:i], nil
			}
			// If we're at EOF, we have a final, non-terminated line. Return it.
			if atEOF {
				return len(data), data, nil
			}
			// Request more data.
			return 0, nil, nil
		},
	)

	result := 0
	for i := 0; i < lines; i++ {
		scanner.Scan()
		result += len(scanner.Bytes()) + len(lineSeparator)
	}
	return result
}

func (p *YamlFileParser) getKeyWithValueOfType(
	node *yaml.Node, element string,
	valueType yaml.Kind, parsedPath *ParsedPath,
) (*yaml.Node, error) {
	for keyIndex := 0; keyIndex < len(node.Content); keyIndex += 2 {
		keyNode := node.Content[keyIndex]
		if keyNode.Kind != yaml.ScalarNode {
			return nil, fmt.Errorf(
				"key node of a map is expected to be a scalar: kind=%v, pos=%v:%v, path=%v",
				keyNode.Kind, keyNode.Line, keyNode.Column, parsedPath.String(),
			)
		}
		if keyNode.Value == element {
			//value goes in the next index after the key in a mapping node
			result := node.Content[keyIndex+1]
			if result.Kind != valueType {
				return nil, fmt.Errorf(
					"middle path node kind is expected to be %v: kind=%v, pos=%v:%v, path=%v", valueType,
					result.Kind, result.Line, result.Column, parsedPath.String(),
				)
			}
			return result, nil
		}
	}
	return nil, fmt.Errorf(
		"failed to find element: pos=%v:%v, path=%v", node.Line, node.Column, parsedPath.String(),
	)
}
