package output

import (
	"bytes"
	"fmt"
	"github.com/aquasecurity/table"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type TableFormatter struct {
	entityName string
}

func (w TableFormatter) FormatMessage(value map[string]EntityInfo) string {
	if len(value) == 0 {
		return ""
	}

	keys := maps.Keys(value)
	slices.Sort(keys)

	buffer := &bytes.Buffer{}

	t := table.New(buffer)
	t.SetDividers(table.UnicodeRoundedDividers)

	var headerRow []string
	headerRow = append(headerRow, w.entityName)
	hasError := anyLineHasError(value)
	if hasError {
		headerRow = append(headerRow, "error")
	}
	entityKeys := getEntityKeys(value)
	headerRow = append(headerRow, entityKeys...)
	t.SetHeaders(headerRow...)

	for _, key := range keys {
		entryValue := value[key]
		var row []string

		row = append(row, key)

		if hasError {
			errorValue := ""
			if entryValue.Error != nil {
				errorValue = entryValue.Error.Error()
			}
			row = append(row, errorValue)
		}

		if entryValue.Result != nil {
			entryValueMap := valueToMap(entryValue.Result)
			for _, entityKey := range entityKeys {
				entryValueValue := entryValueMap[entityKey]
				entryValueValueString := ""
				if entryValueValue != nil {
					entryValueValueString = fmt.Sprintf("%v", entryValueValue)
				}
				row = append(row, entryValueValueString)
			}
		} else {
			for range entityKeys {
				row = append(row, "")
			}
		}

		t.AddRow(row...)
	}

	t.Render()

	return buffer.String()
}

func getEntityKeys(value map[string]EntityInfo) []string {
	for _, entryValue := range value {
		if entryValue.Result != nil {
			return valueKeys(entryValue.Result)
		}
	}
	return nil
}

func anyLineHasError(value map[string]EntityInfo) bool {
	for _, entryValue := range value {
		if entryValue.Error != nil {
			return true
		}
	}
	return false
}
