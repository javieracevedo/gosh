package exec_test

import (
	"gosh/config"
	"gosh/exec"
	"gosh/utils"
	"os"
	"testing"
)

func TestCommandNotExist(t *testing.T) {
    command := []string{"nonexistentcommand"}

    _, err := exec.ExecuteCommand(command);
    if (err == nil) {
        t.Errorf("exec.ExecuteCommand(%v) should return an error when function is called with command that doesn't exists.", command[0])
    }
}

func TestCannotFindCommandEmptyPath(t *testing.T) {
    command := []string{"ls"}

    _, err := exec.ExecuteCommand(command);

    if (err == nil) {
        t.Errorf("exec.ExecuteCommand(%v) should return an error when function is called with command that exists but path is not set.", command[0])
    }
}

func TestExistingCommandPathSet(t *testing.T) {
    command := []string{"ls"}
    initialShellConfig := config.ShellConfig{}
    configLines := []string{
        "path=~/something/:/usr/bin:~/binaries", 
    }
    config.GlobalShellConfig = config.SetConfigValues(initialShellConfig, configLines[0])
    
    _, err := exec.ExecuteCommand(command);

    if (err != nil) {
        t.Errorf("exec.ExecuteCommand(%v) should execute normally and not return an error when existing command is executed with proper path set.", command[0])
    }
}

func TestExistingCommandWithoutReadPermissions(t *testing.T) {
    initialShellConfig := config.ShellConfig{}
    configLines := []string{
        "path=/home/panic/binaries/", 
    }
    config.GlobalShellConfig = config.SetConfigValues(initialShellConfig, configLines[0])

    filename := "/home/panic/binaries/noreadtest"
    permissions := os.FileMode(0333)
    content := "#!/bin/bash \necho 'Hello, World!'"
    utils.CreateFileWithPermissions(filename, permissions, content)

    command := []string{"noreadtest"}
    _, commandExecErr := exec.ExecuteCommand(command);


    if (commandExecErr != nil) {
        t.Errorf("exec.ExecuteCommand(%v) should not return an error when executing a command that exists with execution permission, but not read permissions.", command[0])
    }


    t.Cleanup(func () {
        removeErr := os.Remove(filename);
        if removeErr != nil {
            t.Errorf("Error deleting file: ")
        }
    })
}

func TestExistingCommandWithoutExecutePermissions(t *testing.T) {
    initialShellConfig := config.ShellConfig{}
    configLines := []string{
        "path=/home/panic/binaries/", 
    }
    config.GlobalShellConfig = config.SetConfigValues(initialShellConfig, configLines[0])

    filename := "/home/panic/binaries/noreadtest"
    permissions := os.FileMode(0666)
    content := "#!/bin/bash \necho 'Hello, World!'"
    utils.CreateFileWithPermissions(filename, permissions, content)

    command := []string{"noreadtest"}
    _, commandExecErr := exec.ExecuteCommand(command);


    if (commandExecErr == nil) {
        t.Errorf("exec.ExecuteCommand(%v) should return an error when executing a command that exists without execute permissions", command[0])
    }


    t.Cleanup(func () {
        removeErr := os.Remove(filename);
        if removeErr != nil {
            t.Errorf("Error deleting file: ")
        }
    })
}
