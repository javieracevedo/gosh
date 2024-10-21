package parser_test

import (
	"gosh/parser"
	"gosh/utils"
	"testing"
)


func TestCleanArgsWithEmptyStrings(t *testing.T) {
    expected := []string{"ls", "./", "-la"}
    input := []string{"ls", "", "", "./", "la"}

    output := parser.CleanArgs(input)

    if utils.SlicesEqual(expected, output) {
        t.Errorf("TestCleanArgsWithEmptyStrings: failed")
    }
}

func TestCleanArgsWithTabbedStrings(t *testing.T) {
    expected := []string{"ls", "./", "-la"}
    input := []string{"ls", "\t", "\t", "\t", "\t", "./", "la"}

    output := parser.CleanArgs(input)

    if utils.SlicesEqual(expected, output) {
        t.Errorf("TestCleanArgsWithTabbedStrings: failed")
    }
}

func TestCleanArgsWithNewLinedStrings(t *testing.T) {
    expected := []string{"ls", "./", "-la"}
    input := []string{"ls", "\n\n\n", "\n", "./", "la"}

    output := parser.CleanArgs(input)

    if utils.SlicesEqual(expected, output) {
        t.Errorf("TestCleanArgsWithNewLinedStrings: failed")
    }
}

func TestParseCommandLineWithOnlyNewLinedString(t *testing.T) {
    input := "\n"

	commands, err := parser.ParseCommandLine(input)

    if (commands != nil || err != nil) {
        t.Errorf("TestParseCommandLineWithOnlyNewLinedString: failed")
    }
}

func TestParseCommandLineWithSimpleCommand(t *testing.T) {
    expectedName, expectedArgs := "ls", []string{}
    input := "ls\n"

	commands, _ := parser.ParseCommandLine(input)
	commandName := commands[0][0]
	commandArgs := commands[0][1:]

    if (commandName != expectedName || !utils.SlicesEqual(commandArgs, expectedArgs)) {
        t.Errorf("TestParseCommandLineWithSimpleCommand: failed")
    }
}

func TestParseCommandLineWithCommandAndArgs(t *testing.T) {
    expectedName, expectedArgs := "ls", []string{"./", "-la"}
    input := "ls ./ -la\n"

	commands, _ := parser.ParseCommandLine(input)
	commandName := commands[0][0]
	commandArgs := commands[0][1:]
	
    if (commandName != expectedName || !utils.SlicesEqual(commandArgs, expectedArgs)) {
        t.Errorf("TestParseCommandLineWithCommandAndArgs: failed")
    }
}
