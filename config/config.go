package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const configFilePath string = "gosh.rc"

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

	if (len(cli_args) >= 1) {
		GlobalShellConfig.ShellMode = BATCH
	} else {
		GlobalShellConfig.ShellMode = INTERPRETER
	}

	if strings.HasPrefix(configLine, "qotd_list") {
		GlobalShellConfig.QuoteList = strings.Split(strings.TrimPrefix(configLine, "qotd_list="), ",")
	} else if strings.HasPrefix(configLine, "prompt_color") {
		GlobalShellConfig.PromptColor = strings.ReplaceAll(strings.TrimPrefix(configLine, "prompt_color="), "\\033", "\033")
	} else if strings.HasPrefix(configLine, "path") {
		GlobalShellConfig.Path = strings.TrimPrefix(configLine, "path=")
	}
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

	scanner := bufio.NewScanner(configFile)
	for scanner.Scan() {
		line := scanner.Text()
		setConfigValues(line)
	}
}
