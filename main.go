package main

import (
	"gosh/config"
	"gosh/exec"
)


func main() {
    exec.Clear()

    config.InitShellConfig("./gosh.rc")

    exec.DisplayRandomQuote()
}
