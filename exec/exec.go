package exec

import (
	"fmt"
	"gosh/config"
	"math/rand"
	"os"
	"os/exec"
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
