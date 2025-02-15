package parser

import (
	"bufio"
	"errors"
	"gosh/utils"
	"os"
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
        return nil, nil
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

func ParseBatchFile(fileName string) ([][]string, error) {
    var parsedCommands [][]string

    file, err := os.Open(fileName)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) > 0 {
            lineCommands, err := ParseCommandLine(line)
            if err != nil {
                return nil, err
            }
            if lineCommands != nil {
                parsedCommands = append(parsedCommands, lineCommands...)
            }
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return parsedCommands, nil
}