package output

import (
	"fmt"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/spf13/viper"
	"io"
	"reflect"
	"strings"
)

type Formatter interface {
	FormatMessage(value map[string]EntityInfo) string
}

type EntityInfo struct {
	Result interface{}
	Error  error
}

// Write consumes a map with keys are usually repository names, or any other entity names, like groups
// and values are containers with either structures or errors
func Write(writer io.Writer, entityName string, value map[string]EntityInfo) error {
	formatter, err := createFormatter(entityName)
	if err != nil {
		return err
	}

	message := formatter.FormatMessage(value)

	_, err = fmt.Fprint(writer, message)
	if err != nil {
		return err
	}

	return nil
}

func createFormatter(entityName string) (Formatter, error) {
	outputFormat := config.OutputFormat(viper.GetString("output"))
	if outputFormat == config.JsonOutputFormat {
		return JsonFormatter{entityName}, nil
	}
	if outputFormat == config.LineOutputFormat {
		return LineFormatter{}, nil
	}
	if outputFormat == config.LogOutputFormat {
		return LogFormatter{entityName}, nil
	}
	if outputFormat == config.TableOutputFormat {
		return TableFormatter{entityName}, nil
	}

	return nil, fmt.Errorf("unsupported output format: %v", outputFormat)
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
			entryKeyCamelCase := strings.ToLower(entryKey[:1]) + entryKey[1:]
			entryValue := valueType.Field(i).Interface()
			result[entryKeyCamelCase] = entryValue
		}
	} else {
		if fmt.Sprintf("%v", value) != "" {
			result["result"] = value
		}
	}

	return result
}
