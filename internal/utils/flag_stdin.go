package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

var (
	readFromStdIn bool
	jsonKeyName   string
)

func AddReadFromStdInFlag(command *cobra.Command, jsonKey string) {
	command.Flags().BoolVar(&readFromStdIn, "pipe", false, "Read from stdin after a pipe")
	jsonKeyName = jsonKey
}

func GetReposFromStdInOrDefault(defs []string) ([]string, error) {
	if !readFromStdIn {
		return defs, nil
	}

	if !isInputFromPipe() {
		return nil, errors.New("pipe not found")
	}

	return readFromStdin()
}

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func readFromStdin() ([]string, error) {
	allBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}

	allString := string(allBytes)

	jsonRepos, err := parseReposFromJson(allBytes)
	if err == nil {
		return jsonRepos, nil
	}
	return strings.Fields(allString), nil
}

func parseReposFromJson(value []byte) ([]string, error) {
	var arr []map[string]any
	err := json.Unmarshal(value, &arr)
	if err != nil {
		return nil, err
	}

	var repoNames []string
	for _, r := range arr {
		repoNames = append(repoNames, fmt.Sprintf("%v", r[jsonKeyName]))
	}
	return repoNames, err
}
