
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type LogRotator struct {
    mu           sync.Mutex
    currentFile  *os.File
    filePath     string
    maxSize      int64
    backupCount  int
    currentSize  int64
}

func NewLogRotator(filePath string, maxSizeMB int, backupCount int) (*LogRotator, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024

    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &LogRotator{
        currentFile: file,
        filePath:    filePath,
        maxSize:     maxSize,
        backupCount: backupCount,
        currentSize: stat.Size(),
    }, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.currentSize+int64(len(p)) > lr.maxSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := lr.currentFile.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if err := lr.currentFile.Close(); err != nil {
        return err
    }

    for i := lr.backupCount - 1; i >= 0; i-- {
        oldPath := lr.getBackupPath(i)
        newPath := lr.getBackupPath(i + 1)

        if _, err := os.Stat(oldPath); err == nil {
            if i == lr.backupCount-1 {
                os.Remove(oldPath)
            } else {
                if err := lr.compressAndMove(oldPath, newPath); err != nil {
                    return err
                }
            }
        }
    }

    if err := lr.compressAndMove(lr.filePath, lr.getBackupPath(0)); err != nil {
        return err
    }

    file, err := os.OpenFile(lr.filePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    lr.currentFile = file
    lr.currentSize = 0
    return nil
}

func (lr *LogRotator) compressAndMove(src, dst string) error {
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    dstFile, err := os.Create(dst + ".gz")
    if err != nil {
        return err
    }
    defer dstFile.Close()

    gzWriter := gzip.NewWriter(dstFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, srcFile); err != nil {
        return err
    }

    return os.Remove(src)
}

func (lr *LogRotator) getBackupPath(index int) string {
    if index == 0 {
        return lr.filePath
    }
    ext := filepath.Ext(lr.filePath)
    base := lr.filePath[:len(lr.filePath)-len(ext)]
    return fmt.Sprintf("%s.%d%s", base, index, ext)
}

func (lr *LogRotator) Close() error {
    lr.mu.Lock()
    defer lr.mu.Unlock()
    return lr.currentFile.Close()
}

func main() {
    rotator, err := NewLogRotator("app.log", 10, 5)
    if err != nil {
        panic(err)
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n",
            time.Now().Format(time.RFC3339), i)
        rotator.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}