package output

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type JsonWriter struct {
	entityName string
}

func (w JsonWriter) WriteMessage(value map[string]EntityInfo) string {
	var valueToLog []interface{}

	for key, info := range value {
		entry := createValueToLogEntry(w.entityName, key, info)
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
