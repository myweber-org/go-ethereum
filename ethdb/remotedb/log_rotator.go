
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "sync"
    "time"
)

type RotatingLog struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    currentSize int64
    currentFile *os.File
    fileCount   int
    maxFiles    int
}

func NewRotatingLog(basePath string, maxSize int64, maxFiles int) (*RotatingLog, error) {
    rl := &RotatingLog{
        basePath: basePath,
        maxSize:  maxSize,
        maxFiles: maxFiles,
    }

    if err := rl.openCurrentFile(); err != nil {
        return nil, err
    }

    go rl.cleanupOldFiles()
    return rl, nil
}

func (rl *RotatingLog) openCurrentFile() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentFile != nil {
        rl.currentFile.Close()
    }

    file, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    rl.currentFile = file
    rl.currentSize = info.Size()
    return nil
}

func (rl *RotatingLog) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > rl.maxSize {
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

func (rl *RotatingLog) rotate() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
        rl.currentFile = nil
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)

    if err := os.Rename(rl.basePath, rotatedPath); err != nil {
        return err
    }

    rl.fileCount++

    if err := rl.compressFile(rotatedPath); err != nil {
        return err
    }

    return rl.openCurrentFile()
}

func (rl *RotatingLog) compressFile(sourcePath string) error {
    sourceFile, err := os.Open(sourcePath)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    compressedPath := sourcePath + ".gz"
    destFile, err := os.Create(compressedPath)
    if err != nil {
        return err
    }
    defer destFile.Close()

    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, sourceFile)
    if err != nil {
        return err
    }

    os.Remove(sourcePath)
    return nil
}

func (rl *RotatingLog) cleanupOldFiles() {
    for {
        time.Sleep(1 * time.Hour)

        rl.mu.Lock()
        if rl.fileCount > rl.maxFiles {
            dir := filepath.Dir(rl.basePath)
            baseName := filepath.Base(rl.basePath)

            files, err := filepath.Glob(filepath.Join(dir, baseName+".*.gz"))
            if err == nil && len(files) > rl.maxFiles {
                sortFilesByTimestamp(files)
                filesToRemove := files[:len(files)-rl.maxFiles]
                for _, f := range filesToRemove {
                    os.Remove(f)
                    rl.fileCount--
                }
            }
        }
        rl.mu.Unlock()
    }
}

func sortFilesByTimestamp(files []string) {
    for i := 0; i < len(files); i++ {
        for j := i + 1; j < len(files); j++ {
            ts1 := extractTimestamp(files[i])
            ts2 := extractTimestamp(files[j])
            if ts1 > ts2 {
                files[i], files[j] = files[j], files[i]
            }
        }
    }
}

func extractTimestamp(filename string) string {
    parts := strings.Split(filepath.Base(filename), ".")
    if len(parts) >= 3 {
        return parts[1]
    }
    return ""
}

func (rl *RotatingLog) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    log, err := NewRotatingLog("/var/log/myapp/app.log", 10*1024*1024, 10)
    if err != nil {
        panic(err)
    }
    defer log.Close()

    for i := 0; i < 1000; i++ {
        log.Write([]byte("Log entry " + strconv.Itoa(i) + "\n"))
        time.Sleep(100 * time.Millisecond)
    }
}