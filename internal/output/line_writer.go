package output

import (
	"bytes"
	"fmt"
)

type LineWriter struct {
}

func (w LineWriter) WriteMessage(value map[string]EntityInfo) string {
	buffer := &bytes.Buffer{}

	for key, info := range value {
		infoString := infoToString(info)
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
		valueMap := valueToMap(info.Result)

		for entryKey, entryValue := range valueMap {
			buffer.WriteString(fmt.Sprintf("%s=%s ", entryKey, entryValue))
		}
	}
	return buffer.String()
}
