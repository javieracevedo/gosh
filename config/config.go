package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
)


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


type ShellConfig struct {
    PromptColor string
    Path string
    QuoteList []string
    ShellMode int
}

var GlobalShellConfig ShellConfig;


func SetConfigValues(shellConfig ShellConfig, configLine string) (ShellConfig) {
    cli_args := os.Args[1:]

    trimmedConfigLine := strings.TrimSpace(configLine);


    if (len(cli_args) >= 1) {
        shellConfig.ShellMode = BATCH
    } else {
        shellConfig.ShellMode = INTERPRETER
    }

    if strings.HasPrefix(trimmedConfigLine, "qotd_list") {
        quotes := strings.TrimPrefix(trimmedConfigLine, "qotd_list=")
        if (quotes != "") {
            shellConfig.QuoteList = strings.Split(quotes, ",")
        }
    } else if strings.HasPrefix(trimmedConfigLine, "prompt_color") {
        promptColor := strings.ReplaceAll(strings.TrimPrefix(trimmedConfigLine, "prompt_color="), "\\033", "\033")
        if (promptColor != "") {
            shellConfig.PromptColor = promptColor
        }
    } else if strings.HasPrefix(trimmedConfigLine, "path") {
        shellConfig.Path = strings.TrimPrefix(trimmedConfigLine, "path=")
    }

    return shellConfig
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

func InitShellConfig(configFilePath string) (ShellConfig, error) {
    GlobalShellConfig = ShellConfig{
        PromptColor: DEFAULT_COLOR,
    }

    var configFile, err = os.Open(configFilePath)
    if err != nil {
        if (os.IsTimeout(err)) {
            return GlobalShellConfig, errors.New("timeout error: could not open shell's configuration file")
        } else if (os.IsNotExist(err)) {
            return GlobalShellConfig, errors.New("file not found: could not open shell's configuration file")
        } else if (os.IsPermission(err)) {
            return GlobalShellConfig, errors.New("permission denied: could not open shell's configuration file")
        } else {
            return GlobalShellConfig, errors.New("unknown error: could not open shell's configuration file")
        }
    }
    defer configFile.Close()

    lines, validateFileErr := validateFile(configFile)
    if (validateFileErr != nil) {
        return GlobalShellConfig, validateFileErr
    }

    for _, line := range lines {
        GlobalShellConfig = SetConfigValues(GlobalShellConfig, line)
    }

    return GlobalShellConfig, nil
}

func ExtendPath(directories []string) {
    for _, directory := range directories {
        GlobalShellConfig.Path = GlobalShellConfig.Path + ":" + directory
    }
}
