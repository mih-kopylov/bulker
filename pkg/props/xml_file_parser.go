package props

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"reflect"
)

type XmlFileParser struct {
}

func (j *XmlFileParser) GetProperty(fileBytes []byte, pp *ParsedPath) (*ParsedProperty, error) {
	decoder := xml.NewDecoder(bytes.NewReader(fileBytes))

	//looking for the root xml element
	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, fmt.Errorf("failed to find root node")
			}
			return nil, fmt.Errorf("failed to get next xml token: %w", err)
		}
		//looking for StartElement tokens because they are the ones that describe XML nodes
		if startElementToken, ok := token.(xml.StartElement); ok {
			logrus.WithField("token", startElementToken.Name.Local).Debug("root element found")
			break
		} else {
			continue
		}
	}

	//looking for path inside root element
	for i, element := range pp.Elements {
		for {
			token, err := decoder.Token()
			if err != nil {
				if errors.Is(err, io.EOF) {
					parsedPathNotFound := ParsedPath{pp.Elements[0 : i+1]}
					return nil, fmt.Errorf("failed to find path element: %v", parsedPathNotFound.String())
				}
				return nil, fmt.Errorf("failed to get next xml token: %w", err)
			}

			//looking for StartElement tokens because they are the ones that describe XML nodes
			if startElementToken, ok := token.(xml.StartElement); ok {
				if startElementToken.Name.Local == element {
					break
				}

				err := decoder.Skip()
				if err != nil {
					return nil, fmt.Errorf("failed to skip xml token: %w", err)
				}
			} else {
				continue
			}
		}
	}

	//reading inner text element
	pathToken, err := decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to read next xml token: %w", err)
	}

	if _, ok := pathToken.(xml.EndElement); ok {
		end := int(decoder.InputOffset())
		start := end
		return &ParsedProperty{
			start,
			end,
			"",
			PropertyTypeString,
			QuotesTypeNone,
		}, nil
	} else if charDataToken, ok := pathToken.(xml.CharData); ok {
		end := int(decoder.InputOffset())
		start := end - len(charDataToken)
		return &ParsedProperty{
			start,
			end,
			string(charDataToken),
			PropertyTypeString,
			QuotesTypeNone,
		}, nil
	} else {
		return nil, fmt.Errorf(
			"element expected to be a CharData: path=%v, type=%v", pp.String(), reflect.TypeOf(pathToken),
		)
	}
}
