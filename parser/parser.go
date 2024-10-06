package parser

import (
	"errors"
	"strings"
)

const (
	ARG_MAX = 2097152
)

func CleanArgs(args []string) []string {
    var result []string
    for _, arg := range args {
        if arg != "" && arg != "\n" && arg != "\t" {
            result = append(result, arg)
        }
    }
    return result
}

func ParseCommandLine(line string) ([][]string, error) {
    var commands [][]string;

	trimmedLine := strings.TrimSpace(line)
	if (len(trimmedLine) > ARG_MAX) {
		return nil, errors.New("command exceeds the maximum number of arguments")	
	}

    command := strings.Split(trimmedLine, " ")
	commands = append(commands, command)

    return commands, nil
}
