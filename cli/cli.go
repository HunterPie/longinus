package cli

import (
	"errors"
	"flag"
)

var (
	ErrMissingExecutable = errors.New("missing --executable argument")
)

func verifyRequired(value *string) bool {
	return *value != ""
}

func ParseArguments() (*LonginusArguments, error) {
	executable := flag.String("executable", "", "executable to scan for byte signatures")
	configuration := flag.String("config", "./configuration/default.yaml", "Longinus configuration with byte signature definitions")

	if !verifyRequired(executable) {
		return nil, ErrMissingExecutable
	}

	flag.Parse()

	return &LonginusArguments{
		Executable:    *executable,
		Configuration: *configuration,
	}, nil
}
