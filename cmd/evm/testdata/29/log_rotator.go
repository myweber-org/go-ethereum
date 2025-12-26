package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
    "time"
)

type RotatingLogger struct {
    currentFile   *os.File
    currentSize   int64
    maxFileSize   int64
    basePath      string
    fileCounter   int
    compressOld   bool
}

func NewRotatingLogger(basePath string, maxSizeMB int64, compress bool) (*RotatingLogger, error) {
    if maxSizeMB <= 0 {
        maxSizeMB = 10
    }
    maxBytes := maxSizeMB * 1024 * 1024

    file, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLogger{
        currentFile: file,
        currentSize: info.Size(),
        maxFileSize: maxBytes,
        basePath:    basePath,
        compressOld: compress,
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    if rl.currentSize+int64(len(p)) > rl.maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.currentFile.Write(p)
    if err == nil {
        rl.currentSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.currentFile.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)

    if err := os.Rename(rl.basePath, rotatedPath); err != nil {
        return err
    }

    if rl.compressOld {
        if err := rl.compressFile(rotatedPath); err != nil {
            fmt.Printf("Compression failed: %v\n", err)
        }
    }

    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.currentFile = file
    rl.currentSize = 0
    rl.fileCounter++
    return nil
}

func (rl *RotatingLogger) compressFile(source string) error {
    dest := source + ".gz"
    srcFile, err := os.Open(source)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    destFile, err := os.Create(dest)
    if err != nil {
        return err
    }
    defer destFile.Close()

    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, srcFile); err != nil {
        return err
    }

    if err := os.Remove(source); err != nil {
        return err
    }

    return nil
}

func (rl *RotatingLogger) Close() error {
    return rl.currentFile.Close()
}

func (rl *RotatingLogger) ScanOldFiles() {
    dir := filepath.Dir(rl.basePath)
    baseName := filepath.Base(rl.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return
    }

    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }
        name := entry.Name()
        if strings.HasPrefix(name, baseName+".") && !strings.HasSuffix(name, ".gz") {
            oldPath := filepath.Join(dir, name)
            if rl.compressOld {
                rl.compressFile(oldPath)
            }
        }
    }
}

func main() {
    logger, err := NewRotatingLogger("app.log", 5, true)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    logger.ScanOldFiles()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := logger.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation completed")
}