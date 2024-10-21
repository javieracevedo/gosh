package utils

import "os"

func CreateFileWithPermissions(filename string, perm os.FileMode, content string) error {
    // Create the file with custom permissions
    file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, perm)
    if err != nil {
        return err
    }
    defer file.Close()

    // Write content to the file
    _, err = file.WriteString(content)
    if err != nil {
        return err
    }

    return nil
}