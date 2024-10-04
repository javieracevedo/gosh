package config_test

import (
	"gosh/config"
	"os"
	"strings"
	"testing"
)

func TestSetConfigValuesEmptyFile(t *testing.T) {
    initialShellConfig := config.ShellConfig{}
    configLines := []string{""}
    shellConfig := config.SetConfigValues(initialShellConfig, strings.TrimSpace(configLines[0]))

    if shellConfig.PromptColor != "" || shellConfig.Path != "" || len(shellConfig.QuoteList ) > 0 {
        t.Errorf("config.SetConfigValues(%q), config values should be empty, actual %q", configLines[0], config.GlobalShellConfig)
    }
}

func TestSetConfigValuesWithEmptyValues(t *testing.T) {
    initialShellConfig := config.ShellConfig{}
    configLines := []string{"prompt_color=", "path=", "qotd_list="}
    
    shellConfig := config.SetConfigValues(initialShellConfig, strings.TrimSpace(configLines[0]))
    shellConfig = config.SetConfigValues(shellConfig, strings.TrimSpace(configLines[1]))
    shellConfig = config.SetConfigValues(shellConfig, strings.TrimSpace(configLines[2]))

    if shellConfig.PromptColor != "" || shellConfig.Path != "" || len(shellConfig.QuoteList ) > 0 {
        t.Errorf("config.SetConfigValues(%q | %q | %q), config values should be empty, actual %q", configLines[0], configLines[1], configLines[2], shellConfig)
        return
    }
}

func TestSetConfigValuesWithSetValues(t *testing.T) {
    initialShellConfig := config.ShellConfig{}
    configLines := []string{
        "prompt_color=\033[35m", 
        "path=~/something/:/usr/bin:~/binaries", 
        "qotd_list=All that is gold does not glitter,Not all those who wander are lost,I wish it need not have happened in my time said Frodo.",
    }
    
    shellConfig := config.SetConfigValues(initialShellConfig, configLines[0])
    shellConfig = config.SetConfigValues(shellConfig, configLines[1])
    shellConfig = config.SetConfigValues(shellConfig, configLines[2])

    if shellConfig.PromptColor != "\033[35m" || shellConfig.Path != "~/something/:/usr/bin:~/binaries" || len(shellConfig.QuoteList ) != 3 {
        t.Errorf("config.SetConfigValues(%q | %q | %q), config values should be empty, actual %q", configLines[0], configLines[1], configLines[2], shellConfig)
    }
}

func TestSetConfigValuesWithEmptySpacesInValues(t *testing.T) {
    initialShellConfig := config.ShellConfig{}
    configLines := []string{
        "prompt_color=          ", 
        "path=          ", 
        "qotd_list=          ",
    }
    
    shellConfig := config.SetConfigValues(initialShellConfig, configLines[0])
    shellConfig = config.SetConfigValues(shellConfig, configLines[1])
    shellConfig = config.SetConfigValues(shellConfig, configLines[2])

    if shellConfig.PromptColor != "" || shellConfig.Path != "" || len(shellConfig.QuoteList ) != 0 {
        t.Errorf("config.SetConfigValues(%q | %q | %q), config values should be empty, actual %q", configLines[0], configLines[1], configLines[2], shellConfig)
    }
}

func TestSetConfigValuesWithTabsInValues(t *testing.T) {
    initialShellConfig := config.ShellConfig{}
    configLines := []string{
        "prompt_color=			", 
        "path=				", 
        "qotd_list=				",
    }
    
    shellConfig := config.SetConfigValues(initialShellConfig, configLines[0])
    shellConfig = config.SetConfigValues(shellConfig, configLines[1])
    shellConfig = config.SetConfigValues(shellConfig, configLines[2])

    if shellConfig.PromptColor != "" || shellConfig.Path != "" || len(shellConfig.QuoteList ) != 0 {
        t.Errorf("config.SetConfigValues(%q | %q | %q), config values should be empty, actual %q", configLines[0], configLines[1], configLines[2], shellConfig)
    }
}

func TestInitShellConfigFileNotFound(t *testing.T) {
    filePath := "./non_existent.rc"

    config, initConfigErr := config.InitShellConfig(filePath)

    expectedError := "file not found: could not open shell's configuration file"
    if initConfigErr.Error() != expectedError {
        t.Errorf("config.InitShellConfig(%v), should return (%v) when file doesn't exist", filePath, expectedError)
    }

    if (config.Path == "" && config.PromptColor == "" && len(config.QuoteList) > 0) {
        t.Errorf("config.InitShellConfig(%v), should return empty Path, PromptColor, and QuoteList", filePath)
    }
}

func TestInitShellConfigPermissionDenied(t *testing.T) {
    filePath := "testfile.txt"
    testFile, err := os.Create(filePath)
    if err != nil {
        t.Fatalf("could not create file.")
    }
    testFile.Close()

    err = os.Chmod(filePath, 0220) // 0 + 2 + 0 (-w-)
    if (err != nil) {
        t.Fatalf("could not change file's permissions.")
    }


    config, initConfigErr := config.InitShellConfig(filePath)
    expectedError := "permission denied: could not open shell's configuration file"
    if initConfigErr.Error() != expectedError {
        t.Errorf("config.InitShellConfig(%v), should return (%v) when file's permissions are 0000", filePath, expectedError)
    }

    if config.Path == "" && config.PromptColor == "" && len(config.QuoteList) > 0 {
        t.Errorf("config.InitShellConfig(%v), should return empty Path, PromptColor, and QuoteList", filePath)
    }


    t.Cleanup(func () {
        err = os.Remove(filePath);
        if err != nil {
            t.Errorf("Error deleting file: ")
        }
    })
}    

func TestInitShellConfigInvalidConfigFile(t *testing.T) {
    filePath := "invalid_config_file.txt"
    data := []byte("prompt_color====\npath======\n");

    err := os.WriteFile(filePath, data, 0664)
    if (err != nil) {
        t.Fatalf("could not write to file %v", filePath)
    }


    _, initConfigErr := config.InitShellConfig(filePath)
    expectedErr := "invalid configuration"
    if (initConfigErr.Error() != expectedErr) {
        t.Errorf("config.InitShellConfig(%v), should return 'invalid configuration' error", filePath)
    }


    t.Cleanup(func () {
        // os cleanup?
        err = os.Remove(filePath);
        if err != nil {
            t.Errorf("Error deleting file: ")
        }
    })
}

func TestInitShellConfigValidConfigFile(t *testing.T) {
    filePath := "valid_file.rc"
    data := []byte("prompt_color=\033[35m\npath=/usr/bin\nqotd_list=quote 1,quote 2");

    err := os.WriteFile(filePath, data, 0664)
    if (err != nil) {
        t.Fatalf("could not write to file %v", filePath)
    }


    config, initConfigErr := config.InitShellConfig(filePath)
    if (initConfigErr != nil) {
        t.Errorf("config.InitShellConfig(%v), should return a valid configuration file and not error", filePath)
    }

    if (config.Path != "/usr/bin" || config.PromptColor != "\033[35m" || config.QuoteList[0] != "quote 1" || config.QuoteList[1] != "quote 2") {
        t.Errorf("config path doesnt have right values")
    }


    t.Cleanup(func () {
        err = os.Remove(filePath);
        if err != nil {
            t.Errorf("Error deleting file: ")
        }
    })
}
