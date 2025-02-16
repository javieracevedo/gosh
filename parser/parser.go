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

type ParsedCommand struct {
    Name string
    Argv []string
    RedirectFilePath string
}

func ParseCommandLine(line string) ([]ParsedCommand, error) {
    var commands []ParsedCommand;

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
        
        // See if there is a redirect line (eg: command >> file.txt)
        // if so, parse it.
        splittedByRedirectCommandLine := strings.Split(strings.TrimSpace(splittedCommands[i]), ">>")
        var redirectFilePath string
        if (len(splittedByRedirectCommandLine) > 1) {
            redirectFilePath = strings.TrimSpace(splittedByRedirectCommandLine[1])
        }

        command := strings.Split(strings.TrimSpace(splittedByRedirectCommandLine[0]), " ")
        cleanedCommand := utils.RemoveEmptyStrings(command)
        parsedCommand := ParsedCommand{
            Name: command[0],
            Argv: command,
            RedirectFilePath: redirectFilePath,
        }

        if (len(cleanedCommand) >= 1) {
            commands = append(commands, parsedCommand)
        }
    }

    return commands, nil
}

func ParseBatchFile(fileName string) ([]ParsedCommand, error) {
    var parsedCommands []ParsedCommand

    file, err := os.Open(fileName)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) > 0 {
            parsedCommandLine, err := ParseCommandLine(line)
            if err != nil {
                return nil, err
            }

            for i := 0; i < len(parsedCommandLine); i++ {
                parsedCommands = append(parsedCommands, parsedCommandLine[i])
            }
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return parsedCommands, nil
}