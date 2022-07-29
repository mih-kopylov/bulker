package output

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"reflect"
)

// LogWriter Writes the information with a regular logging
type LogWriter struct {
	entityName string
}

func (w LogWriter) WriteMessage(value map[string]EntityInfo) string {
	buffer := &bytes.Buffer{}

	keys := maps.Keys(value)
	slices.Sort(keys)

	for _, key := range keys {
		logger := logrus.New()
		logger.SetOutput(buffer)
		loggerEntry := logger.WithField(w.entityName, key)

		info := value[key]
		loggerEntry = addLoggerEntries(loggerEntry, info.Result)
		if info.Error != nil {
			loggerEntry.WithError(info.Error).Errorln()
		} else {
			loggerEntry.Infoln()
		}
	}
	return buffer.String()
}

func addLoggerEntries(entry *logrus.Entry, value interface{}) *logrus.Entry {
	valueMap := valueToMap(value)

	if valueMap == nil {
		return entry
	}

	for entryKey, entryValue := range valueMap {
		entry = entry.WithField(entryKey, entryValue)
	}

	return entry
}

func valueToMap(value interface{}) map[string]interface{} {
	if value == nil {
		return nil
	}

	result := map[string]interface{}{}

	valueType := reflect.Indirect(reflect.ValueOf(value))
	if valueType.Kind() == reflect.Map {
		for _, mapKey := range valueType.MapKeys() {
			entryKey := fmt.Sprintf("%v", mapKey.Interface())
			mapValue := valueType.MapIndex(mapKey)
			entryValue := mapValue.Interface()
			result[entryKey] = entryValue
		}
	} else if valueType.Kind() == reflect.Struct {
		for i := 0; i < valueType.NumField(); i++ {
			entryKey := valueType.Type().Field(i).Name
			entryValue := valueType.Field(i).Interface()
			result[entryKey] = entryValue
		}
	} else {
		if fmt.Sprintf("%v", value) != "" {
			result["result"] = value
		}
	}

	return result
}
