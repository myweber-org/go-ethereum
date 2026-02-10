package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
)

const maxFileSize = 1024 * 1024 // 1MB

func rotateLogFile(filePath string) error {
    fileInfo, err := os.Stat(filePath)
    if err != nil {
        return err
    }

    if fileInfo.Size() < maxFileSize {
        return nil
    }

    for i := 1; ; i++ {
        newPath := filePath + "." + strconv.Itoa(i)
        if _, err := os.Stat(newPath); os.IsNotExist(err) {
            err := os.Rename(filePath, newPath)
            if err != nil {
                return err
            }
            fmt.Printf("Rotated log file to: %s\n", newPath)
            break
        }
    }

    newFile, err := os.Create(filePath)
    if err != nil {
        return err
    }
    newFile.Close()
    return nil
}

func writeLog(filePath, message string) error {
    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = io.WriteString(file, message+"\n")
    if err != nil {
        return err
    }

    return rotateLogFile(filePath)
}

func main() {
    logPath := filepath.Join(".", "application.log")
    for i := 0; i < 15000; i++ {
        err := writeLog(logPath, fmt.Sprintf("Log entry number: %d", i))
        if err != nil {
            fmt.Printf("Error writing log: %v\n", err)
            break
        }
    }
}