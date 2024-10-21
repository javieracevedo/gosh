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

func RemoveEmptyStrings(input []string) []string {
    var result []string
    for _, str := range input {
        if str != "" {
            result = append(result, str)
        }
    }
    return result
}

func ConvertTabToSpaces(input string) string {
    var result string
    for _, char := range input {
        if char == '\t' {
            result += string(' ')
        } else {
			result = result + string(char)
		}
    }
    return result
}
