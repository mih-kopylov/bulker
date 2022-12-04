package config

import "fmt"

type OutputFormat string

func (g *OutputFormat) String() string {
	return string(*g)
}

func (g *OutputFormat) Set(v string) error {
	switch v {
	case string(JsonOutputFormat), string(LineOutputFormat), string(LogOutputFormat), string(TableOutputFormat):
		*g = OutputFormat(v)
		return nil
	default:
		return fmt.Errorf(
			"must be one of '%s' '%s' '%s' '%s'", JsonOutputFormat, LineOutputFormat, LogOutputFormat,
			TableOutputFormat,
		)
	}
}

func (g *OutputFormat) Type() string {
	return "OutputFormat"
}

const (
	JsonOutputFormat  OutputFormat = "json"
	LineOutputFormat  OutputFormat = "line"
	LogOutputFormat   OutputFormat = "log"
	TableOutputFormat OutputFormat = "table"
)
