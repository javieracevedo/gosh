package main

import (
	"fmt"
	"gosh/config"
	"log"
	"math/rand"
	"os"
	"os/exec"
)


func exec_clear() {
    clear_cmd := exec.Command("clear")
    clear_cmd.Stdout = os.Stdout
    clear_cmd.Run()
}

func printRandomQuote() {
    quote_length := len(config.GlobalShellConfig.QuoteList)
    if (quote_length > 0) {
        randomIndex := rand.Intn(len(config.GlobalShellConfig.QuoteList))
        randomItem := config.GlobalShellConfig.QuoteList[randomIndex]
        fmt.Print(config.CYAN_COLOR, "Quote of the Day: \n\n", config.DEFAULT_COLOR)
        fmt.Print(config.CYAN_COLOR, randomItem, "\n\n", config.DEFAULT_COLOR)
    }
}

func main() {
    exec_clear()

    _, initShellConfigErr := config.InitShellConfig("./gosh.rc")
    if (initShellConfigErr != nil) {
        log.Fatal(initShellConfigErr)
    }
    printRandomQuote()
}
