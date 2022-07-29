package output

import (
	"fmt"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/spf13/viper"
)

type Writer interface {
	WriteMessage(value map[string]EntityInfo) string
}

type EntityInfo struct {
	Result interface{}
	Error  error
}

// Write consumes a map with keys are usually repository names, or any other entity names, like groups
// and values are containers with either structures or errors
func Write(entityName string, value map[string]EntityInfo) error {
	writer, err := createWriter(entityName)
	if err != nil {
		return err
	}

	message := writer.WriteMessage(value)
	fmt.Print(message)

	return nil
}

func createWriter(entityName string) (Writer, error) {
	outputFormat := config.OutputFormat(viper.GetString("output"))
	if outputFormat == config.JsonOutputFormat {
		return JsonWriter{entityName}, nil
	}
	if outputFormat == config.LineOutputFormat {
		return LineWriter{}, nil
	}
	if outputFormat == config.LogOutputFormat {
		return LogWriter{entityName}, nil
	}
	if outputFormat == config.TableOutputFormat {
		return TableWriter{entityName}, nil
	}

	return nil, fmt.Errorf("unsupported output format: %v", outputFormat)
}
