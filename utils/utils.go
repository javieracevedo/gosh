package utils

import (
	"fmt"
	"gosh/config"
)

func GoshPrint(val any) {
    fmt.Print(config.GlobalShellConfig.PromptColor, val, "\n", config.DEFAULT_COLOR)
}

func SlicesEqual(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for idx := range a {
		if a[idx] != b[idx] {
			return false
		}
	}
	
	return true
}
