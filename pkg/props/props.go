package props

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"strings"
)

func GetPropertyFromFile(fileName string, propertyPath string) (string, error) {
	pp, err := ParsePath(propertyPath)
	if err != nil {
		return "", fmt.Errorf("failed to parse path expression: %w", err)
	}

	parser, err := getFileParser(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create file parser: %w", err)
	}

	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	prop, err := parser.GetProperty(fileBytes, pp)
	if err != nil {
		return "", errors.Wrap(
			ErrPropertyNotFound, fmt.Sprintf(
				"failed to get property: fileName=%v path=%v %v", fileName,
				propertyPath, err.Error(),
			),
		)
	}

	return prop.Text, nil
}

func getFileParser(fileName string) (FileParser, error) {
	if strings.HasSuffix(fileName, ".json") {
		return &JsonFileParser{}, nil
	}
	if strings.HasSuffix(fileName, ".yaml") || strings.HasSuffix(fileName, ".yml") {
		return &YamlFileParser{}, nil
	}
	if strings.HasSuffix(fileName, ".xml") {
		return &XmlFileParser{}, nil
	}
	return nil, fmt.Errorf("filename=%v: %w", fileName, ErrUnsupportedFileFormat)
}

func ParsePath(propertyPath string) (*ParsedPath, error) {
	if !strings.HasPrefix(propertyPath, "$.") {
		return nil, fmt.Errorf("path=%v: %w", propertyPath, ErrMalformedPathExpressoin)
	}
	propertyPath = propertyPath[2:]
	elements := strings.Split(propertyPath, ".")
	return &ParsedPath{Elements: elements}, nil
}
