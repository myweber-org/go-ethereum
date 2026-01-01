package main

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
)

func rotateLog(logPath string) error {
    info, err := os.Stat(logPath)
    if err != nil {
        if os.IsNotExist(err) {
            return nil
        }
        return err
    }

    if info.Size() == 0 {
        return nil
    }

    timestamp := time.Now().Format("20060102_150405")
    ext := filepath.Ext(logPath)
    base := logPath[:len(logPath)-len(ext)]
    archivePath := fmt.Sprintf("%s_%s%s", base, timestamp, ext)

    err = os.Rename(logPath, archivePath)
    if err != nil {
        return err
    }

    file, err := os.Create(logPath)
    if err != nil {
        return err
    }
    return file.Close()
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "Usage: %s <logfile>\n", os.Args[0])
        os.Exit(1)
    }

    logPath := os.Args[1]
    err := rotateLog(logPath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error rotating log: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Log rotated successfully: %s\n", logPath)
}