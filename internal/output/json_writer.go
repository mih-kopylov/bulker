package output

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type JsonWriter struct {
	entityName string
}

func (w JsonWriter) WriteMessage(value map[string]EntityInfo) string {
	var valueToLog []interface{}

	keys := maps.Keys(value)
	slices.Sort(keys)

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
