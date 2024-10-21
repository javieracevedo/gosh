package main

import (
	"bufio"
	"fmt"
	"gosh/config"
	"gosh/exec"
	"gosh/parser"
	"log"
	"os"
	"syscall"
)


func main() {
    config.InitShellConfig("./gosh.rc")
    
    var ws syscall.WaitStatus
    pid, _ := exec.ExecuteCommand([]string{"clear"})
    syscall.Wait4(pid, &ws, 0, nil)

    exec.DisplayRandomQuote()

    reader := bufio.NewReader(os.Stdin)

    for {
        fmt.Print("gosh> ")
    
        input, _ := reader.ReadString('\n')
    
        commands, err := parser.ParseCommandLine(input)
        if err != nil {
            log.Fatal(err)
        }

        exec.ExecuteCommandsAndWait(commands)
    }
}
