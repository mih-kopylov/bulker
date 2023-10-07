package output

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// LogFormatter Formats the information with a regular logging
type LogFormatter struct {
	entityName string
}

func (w LogFormatter) FormatMessage(value map[string]EntityInfo) string {
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
	keys := valueKeys(value)
	valueMap := valueToMap(value)

	if valueMap == nil {
		return entry
	}

	for _, key := range keys {
		entry = entry.WithField(key, valueMap[key])
	}

	return entry
}
