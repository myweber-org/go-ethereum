package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type RotatingLogger struct {
    currentFile *os.File
    filePath    string
    bytesWritten int64
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
    file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
        filePath:    path,
        bytesWritten: info.Size(),
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    if rl.bytesWritten+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.currentFile.Write(p)
    rl.bytesWritten += int64(n)
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.currentFile.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", rl.filePath, timestamp)
    if err := os.Rename(rl.filePath, rotatedPath); err != nil {
        return err
    }

    if err := rl.compressFile(rotatedPath); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.currentFile = file
    rl.bytesWritten = 0
    rl.cleanupOldBackups()
    return nil
}

func (rl *RotatingLogger) compressFile(path string) error {
    src, err := os.Open(path)
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(path + ".gz")
    if err != nil {
        return err
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    if _, err := io.Copy(gz, src); err != nil {
        return err
    }

    os.Remove(path)
    return nil
}

func (rl *RotatingLogger) cleanupOldBackups() {
    pattern := rl.filePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackups {
        return
    }

    for i := 0; i < len(matches)-maxBackups; i++ {
        os.Remove(matches[i])
    }
}

func (rl *RotatingLogger) Close() error {
    return rl.currentFile.Close()
}

func main() {
    logger, err := NewRotatingLogger("app.log")
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }
}
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type LogRotator struct {
    currentFile *os.File
    currentSize int64
    basePath    string
    fileIndex   int
}

func NewLogRotator(basePath string) (*LogRotator, error) {
    rotator := &LogRotator{
        basePath: basePath,
    }

    if err := rotator.openCurrentFile(); err != nil {
        return nil, err
    }

    return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
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
    if lr.currentFile != nil {
        lr.currentFile.Close()
        if err := lr.compressFile(lr.currentFile.Name()); err != nil {
            return err
        }
    }

    lr.fileIndex++
    return lr.openCurrentFile()
}

func (lr *LogRotator) openCurrentFile() error {
    filename := fmt.Sprintf("%s.%d.log", lr.basePath, lr.fileIndex)
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    lr.currentFile = file
    lr.currentSize = stat.Size()

    lr.cleanupOldBackups()
    return nil
}

func (lr *LogRotator) compressFile(source string) error {
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

    _, err = io.Copy(gzWriter, srcFile)
    if err != nil {
        return err
    }

    return os.Remove(source)
}

func (lr *LogRotator) cleanupOldBackups() {
    files, err := filepath.Glob(lr.basePath + ".*.log.gz")
    if err != nil {
        return
    }

    if len(files) <= maxBackups {
        return
    }

    var fileInfos []struct {
        path string
        time time.Time
    }

    for _, file := range files {
        parts := strings.Split(file, ".")
        if len(parts) < 3 {
            continue
        }

        idx, err := strconv.Atoi(parts[len(parts)-3])
        if err != nil {
            continue
        }

        info, err := os.Stat(file)
        if err != nil {
            continue
        }

        fileInfos = append(fileInfos, struct {
            path string
            time time.Time
        }{file, info.ModTime()})
    }

    sortFilesByTime(fileInfos)

    for i := 0; i < len(fileInfos)-maxBackups; i++ {
        os.Remove(fileInfos[i].path)
    }
}

func sortFilesByTime(files []struct {
    path string
    time time.Time
}) {
    for i := 0; i < len(files); i++ {
        for j := i + 1; j < len(files); j++ {
            if files[i].time.After(files[j].time) {
                files[i], files[j] = files[j], files[i]
            }
        }
    }
}

func (lr *LogRotator) Close() error {
    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("application")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }

        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}