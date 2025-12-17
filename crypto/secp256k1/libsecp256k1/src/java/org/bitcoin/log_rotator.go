package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
)

const (
    maxFileSize  = 1024 * 1024 // 1MB
    maxBackups   = 5
    logFileName  = "app.log"
)

func rotateLogFile(filePath string) error {
    for i := maxBackups - 1; i > 0; i-- {
        oldPath := filePath + "." + strconv.Itoa(i)
        newPath := filePath + "." + strconv.Itoa(i+1)
        if _, err := os.Stat(oldPath); err == nil {
            os.Rename(oldPath, newPath)
        }
    }

    if _, err := os.Stat(filePath); err == nil {
        backupPath := filePath + ".1"
        os.Rename(filePath, backupPath)
    }
    return nil
}

func checkLogSize(filePath string) (bool, error) {
    fileInfo, err := os.Stat(filePath)
    if err != nil {
        if os.IsNotExist(err) {
            return false, nil
        }
        return false, err
    }
    return fileInfo.Size() >= maxFileSize, nil
}

func writeLog(message string) error {
    filePath := logFileName
    needsRotation, err := checkLogSize(filePath)
    if err != nil {
        return err
    }

    if needsRotation {
        err = rotateLogFile(filePath)
        if err != nil {
            return err
        }
    }

    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = io.WriteString(file, message+"\n")
    return err
}

func cleanupOldBackups() {
    for i := maxBackups + 1; i <= maxBackups+10; i++ {
        backupPath := logFileName + "." + strconv.Itoa(i)
        os.Remove(backupPath)
    }
}

func main() {
    cleanupOldBackups()

    for i := 0; i < 100; i++ {
        logMessage := fmt.Sprintf("Log entry %d: Application event occurred", i)
        err := writeLog(logMessage)
        if err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
        }
    }

    fmt.Println("Log rotation test completed")
    fmt.Printf("Current log files in %s:\n", filepath.Dir(logFileName))
    
    files, _ := filepath.Glob(logFileName + "*")
    for _, file := range files {
        if strings.HasPrefix(filepath.Base(file), logFileName) {
            fmt.Println("  ", file)
        }
    }
}