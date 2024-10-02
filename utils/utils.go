package utils

import (
	"fmt"
	"gosh/config"
)

func GoshPrint(val any) {
	fmt.Print(config.GlobalShellConfig.PromptColor, val, "\n", config.DEFAULT_COLOR)
}