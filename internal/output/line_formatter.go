package output

import (
	"bytes"
	"fmt"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"strings"
)

type LineFormatter struct {
}

func (w LineFormatter) FormatMessage(value map[string]EntityInfo) string {
	buffer := &bytes.Buffer{}

	keys := maps.Keys(value)
	slices.Sort(keys)

	for _, key := range keys {
		infoString := infoToString(value[key])
		if infoString == "" {
			buffer.WriteString(fmt.Sprintln(key))
		} else {
			buffer.WriteString(fmt.Sprintf("%v: %v\n", key, infoString))
		}
	}

	return buffer.String()
}

func infoToString(info EntityInfo) string {
	buffer := &bytes.Buffer{}
	if info.Error != nil {
		buffer.WriteString(info.Error.Error())
		buffer.WriteString(" ")
	}
	if info.Result != nil {
		keys := valueKeys(info.Result)
		valueMap := valueToMap(info.Result)

		for _, key := range keys {
			buffer.WriteString(fmt.Sprintf("%s=%s ", key, valueMap[key]))
		}
	}
	return strings.TrimSpace(buffer.String())
}
