package exec_test

import (
	"errors"
	"gosh/config"
	"gosh/exec"
	"gosh/utils"
	"os"
	"syscall"
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
    _, commandExecErr := exec.ExecuteCommand(command)
    
    if commandExecErr == nil {
        t.Errorf("exec.ExecuteCommand(%v) should return an error when executing a shell script without read permissions.", command[0])
    }
    
    t.Cleanup(func() {
        removeErr := os.Remove(filename)
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

func TestBuiltInCd_EACCESS(t *testing.T) {
    dirname := "/tmp/no_access_dir"
    err := os.Mkdir(dirname, 0600)
    if err != nil {
        t.Fatalf("Failed to create test directory: %v", err)
    }

    cdError := exec.BuiltinCd(dirname)

    if cdError == nil {
        t.Errorf("BuiltinCd(%v) should return an error when changing to directory without execute permissions", dirname)
    }
    if !errors.Is(cdError, syscall.EACCES) {
        t.Errorf("Expected EACCES error, got: %v", cdError)
    }

    t.Cleanup(func() {
        os.Remove(dirname)
    })
}

func TestBuiltInCd_ENAMETOOLONG(t *testing.T) {
	dirname := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	
	cdError := exec.BuiltinCd(dirname)

	if cdError == nil {
		t.Errorf("BuiltInCd(%v) should return an error when changing to directory with name longer than 255 characters.", dirname)
	}

	if !errors.Is(cdError, syscall.ENAMETOOLONG) {
		t.Errorf("Expected ENAMETOOLONG error, got: %v", cdError)
	}
}

func TestBuiltInCd_ENOENT(t *testing.T) {
	dirname := "nonexistent_file"

	cdError := exec.BuiltinCd(dirname)

	if cdError == nil {
		t.Errorf("BuiltInCd(%v) should return an error", dirname)
	}

	if !errors.Is(cdError, syscall.ENOENT) {
		t.Errorf("Expected ENOENT error, got: %v", cdError)
	}
}

func TestBuiltInCd_ENOTDIR(t *testing.T) {
	fileName := "text.txt"
	_, err := os.Create(fileName)
	if (err != nil) {
		t.Fatalf("Failed to create test file (text.txt): %v", err)
	}

	cdError := exec.BuiltinCd(fileName)
	if (cdError == nil) {
		t.Errorf("BuiltInCd(%v) should return an error when attempting to cd into a regular file", fileName)
	}

	if !errors.Is(cdError, syscall.ENOTDIR) {
		t.Errorf("Expected ENOTDIR error, got: %v", cdError)
	}

	t.Cleanup(func() {
        os.Remove(fileName)
    })
}

func TestBuiltInCd_ELOOP(t *testing.T) {
	link1 := "./test_files/link1"

	cdError := exec.BuiltinCd(link1)

	if (cdError == nil) {
		t.Errorf("BuiltInCd(%v) should return an error when attempting to cd into circular symbolic links", link1)
	}

	if !errors.Is(cdError, syscall.ELOOP) {
		t.Errorf("Expected EIO error, got: %v", cdError)
	}
}