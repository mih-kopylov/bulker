package output

import (
	"encoding/json"
	"maps"
	"slices"

	"github.com/sirupsen/logrus"
)

type JsonFormatter struct {
	entityName string
}

func (w JsonFormatter) FormatMessage(value map[string]EntityInfo) string {
	// An empty slice is created in order to have an empty json when the slice is marshalled without any values
	//goland:noinspection GoPreferNilSlice
	valueToLog := []any{}

	keys := slices.Sorted(maps.Keys(value))

	for _, key := range keys {
		entry := createValueToLogEntry(w.entityName, key, value[key])
		valueToLog = append(valueToLog, entry)
	}

	result, err := json.Marshal(valueToLog)
	if err != nil {
		logrus.Panicf("failed to marshal a map to json: map=%v err=%v", valueToLog, err)
	}
	return string(result)
}

func createValueToLogEntry(entityName string, key string, info EntityInfo) map[string]interface{} {
	result := map[string]interface{}{}

	result[entityName] = key

	if info.Error != nil {
		result["error"] = info.Error.Error()
	}

	valueMap := valueToMap(info.Result)
	for entryKey, entryValue := range valueMap {
		result[entryKey] = entryValue
	}

	return result
}
