package config

import "fmt"

// RunMode implements Value in spf13/pflag for custom flag type
type RunMode string

func (rm *RunMode) String() string {
	return string(*rm)
}

func (rm *RunMode) Set(v string) error {
	switch v {
	case string(Parallel), string(Sequential):
		*rm = RunMode(v)
		return nil
	default:
		return fmt.Errorf("must be either '%s' or '%s'", Parallel, Sequential)
	}
}

func (rm *RunMode) Type() string {
	return "RunMode"
}

const (
	Parallel   RunMode = "par"
	Sequential RunMode = "seq"
)
