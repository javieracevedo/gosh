package exec

import (
    "errors"
    "fmt"
    "gosh/config"
    "math/rand"
    "os"
    "os/exec"
    "strings"
    "sync"
    "syscall"
    "gosh/parser"

    "golang.org/x/sys/unix"
)


func Clear() {
    clear_cmd := exec.Command("clear")
    clear_cmd.Stdout = os.Stdout
    clear_cmd.Run()
}

func DisplayRandomQuote() {
    quote_length := len(config.GlobalShellConfig.QuoteList)
    if (quote_length > 0) {
        randomIndex := rand.Intn(len(config.GlobalShellConfig.QuoteList))
        randomItem := config.GlobalShellConfig.QuoteList[randomIndex]
        fmt.Print(config.CYAN_COLOR, "Quote of the Day: \n\n", config.DEFAULT_COLOR)
        fmt.Print(config.CYAN_COLOR, randomItem, "\n\n", config.DEFAULT_COLOR)
    }
}

func GetCommandPath(commandName string) (string, error) {
    splitPath := strings.Split(config.GlobalShellConfig.Path, ":")
    var commandPath string;

    for _, path := range splitPath {
        if (strings.HasSuffix(path, "/")) {
            commandPath = path + commandName;
        } else {
            commandPath = path + "/" + commandName
        }

        err := unix.Access(commandPath, unix.R_OK)
        if (err == nil) {
            return commandPath, nil
        }
    }
    return "", errors.New("could not find command in shell's path")
}

func ExecuteCommand(argv []string, stdoutFilePath string) (int, error) {
    if len(argv) < 1 {
        return -1, errors.New("empty command")
    }

    commandPath, getCommandPathErr := GetCommandPath(argv[0])
    if getCommandPathErr != nil {
        return -1, errors.New("could not find command " + argv[0])
    }
    
    currentWorkDirectory, getcwderr := os.Getwd()
    if getcwderr != nil {
        return -1, errors.New("could not get current working directory")
    }
    
    path := fmt.Sprintf("PATH=%s", config.GlobalShellConfig.Path)

    // Default file descriptors
    files := []uintptr{0, 1, 2} // stdin, stdout, stderr

    // If stdout file path is provided, open/create the file and use its descriptor
    if stdoutFilePath != "" {
        stdoutFile, err := os.OpenFile(stdoutFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
        if err != nil {
            return -1, fmt.Errorf("error opening output file: %w", err)
        }
        defer stdoutFile.Close()
        files[1] = stdoutFile.Fd()
    }

    attr := &syscall.ProcAttr{
        Dir: currentWorkDirectory,
        Env: append(os.Environ(), path),
        Files: files,
    }
    
    pid, err := syscall.ForkExec(commandPath, argv, attr)
    if err != nil {
        return -1, errors.New("error while forking and executing command")
    }

    return pid, nil
}

func BuiltinCd(directory string) error {
    if err := syscall.Chdir(directory); err != nil {
        error := fmt.Errorf("error changing directory: %w", err)
        return error
    }
    return nil;
}

func BuiltinPath(command parser.ParsedCommand) {
    if (len(command.Argv) > 1) {
        config.ExtendPath(command.Argv[1:])
    } else {
        fmt.Println(config.GlobalShellConfig.Path)
    }
}

func ExecuteCommands(wg *sync.WaitGroup, commands []parser.ParsedCommand) {
    for i := 0; i < len(commands); i++ {
        wg.Add(1)
        go func(i int, command parser.ParsedCommand) {
            defer wg.Done()
            if command.Name == "cd" {
                if len(command.Argv) > 1 {
                    cdError := BuiltinCd(command.Argv[1])
                    if (cdError != nil) {
                        fmt.Println(cdError)
                    }
                }
            } else if command.Name == "exit" {
                syscall.Exit(0)
            } else if (command.Name == "path")  {
                BuiltinPath(command)
            } else {
                pid, err := ExecuteCommand(command.Argv, command.RedirectFilePath)
                if err != nil {
                    fmt.Println("error executing command:", err)
                    return
                }
                var wstatus syscall.WaitStatus
                if _, err := syscall.Wait4(pid, &wstatus, 0, nil); err != nil {
                    fmt.Println("error waiting for child process:", err)
                    return
                }
            }
        }(i, commands[i])
    }
}

func ExecuteCommandsAndWait(commands []parser.ParsedCommand) {
    var wg sync.WaitGroup

    ExecuteCommands(&wg, commands)

    wg.Wait()
}
