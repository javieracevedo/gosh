package config

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

const configFilePath string = "./gosh.rc"

const (
	BATCH = iota
	INTERPRETER
)

const (
	RED_COLOR     = "\033[31m"
	GREEN_COLOR   = "\033[32m"
	YELLOW_COLOR  = "\033[33m"
	BLUE_COLOR    = "\033[34m"
	MAGENTA_COLOR = "\033[35m"
	CYAN_COLOR    = "\033[36m"
	DEFAULT_COLOR   = "\033[0m"
)

const (
	F_OK = 0
	X_OK = 1
	W_OK = 2 
	R_OK = 4
)


type shellConfig struct {
    PromptColor string
    Path string
	QuoteList []string
	ShellMode int
}

var GlobalShellConfig shellConfig;


func setConfigValues(configLine string) {
	cli_args := os.Args[1:]

	trimmedConfigLine := strings.TrimSpace(configLine);


	if (len(cli_args) >= 1) {
		GlobalShellConfig.ShellMode = BATCH
	} else {
		GlobalShellConfig.ShellMode = INTERPRETER
	}

	if strings.HasPrefix(trimmedConfigLine, "qotd_list") {
		quotes := strings.TrimPrefix(trimmedConfigLine, "qotd_list=")
		if (quotes != "") {
			GlobalShellConfig.QuoteList = strings.Split(quotes, ",")
		}
	} else if strings.HasPrefix(trimmedConfigLine, "prompt_color") {
		promptColor := strings.ReplaceAll(strings.TrimPrefix(trimmedConfigLine, "prompt_color="), "\\033", "\033")
		if (promptColor != "") {
			GlobalShellConfig.PromptColor = strings.ReplaceAll(strings.TrimPrefix(trimmedConfigLine, "prompt_color="), "\\033", "\033")
		}
	} else if strings.HasPrefix(trimmedConfigLine, "path") {
		GlobalShellConfig.Path = strings.TrimPrefix(trimmedConfigLine, "path=")
	}
}

func validateFile(configFile *os.File) ([]string, error) {
	scanner := bufio.NewScanner(configFile)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		
		if (line == "") {
			continue
		}

		delimCount := strings.Count(line, "=")
		if (delimCount != 1) {
			return nil, errors.New("invalid configuration")
		}

		lines = append(lines, line)
	}
	return lines, nil
}

func InitShellConfig() {
	GlobalShellConfig = shellConfig{
		PromptColor: DEFAULT_COLOR,
	}

	var configFile, err = os.Open(configFilePath)
	if err != nil {
		errorMessage := "could not find shell's configuration file (gshell.rc), using default values..."
		if (os.IsTimeout(err)) {
			fmt.Println("Timeout:", errorMessage)
		} else if (os.IsNotExist(err)) {
			fmt.Println("File not found:", errorMessage)
		} else if (os.IsPermission(err)) {
			fmt.Println("Permission denied:", errorMessage)
		} else {
			fmt.Println("Unknown error:", errorMessage)
		}
	}
	defer configFile.Close()

	lines, validateFileErr := validateFile(configFile)
	if (validateFileErr != nil) {
		log.Fatal(validateFileErr)
	}

	for _, line := range lines {
		setConfigValues(strings.TrimSpace(line))
	}
}
