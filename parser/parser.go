package parser

import (
	"errors"
	"gosh/utils"
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

    if line == "\n" {
        return commands, nil
    }

	line = utils.ConvertTabToSpaces(line)

    if (len(line) > ARG_MAX) {
        return nil, errors.New("command exceeds the maximum number of arguments")	
    }

    splittedCommands := strings.Split(strings.TrimSpace(line), "&")
    for i := 0; i < len(splittedCommands); i++ {
        trimmedCommand := strings.TrimSpace(splittedCommands[i])
        if (trimmedCommand == "") {
            continue
        }

        command := strings.Split(strings.TrimSpace(splittedCommands[i]), " ")
		cleanedCommand := utils.RemoveEmptyStrings(command)

        if (len(cleanedCommand) >= 1) {
            commands = append(commands, cleanedCommand)
        }
    }

    return commands, nil
}
