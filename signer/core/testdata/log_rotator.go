package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

const (
    maxFileSize    = 10 * 1024 * 1024 // 10MB
    maxBackupFiles = 5
    logFileName    = "app.log"
)

type LogRotator struct {
    currentFile *os.File
    currentSize int64
    basePath    string
}

func NewLogRotator(logDir string) (*LogRotator, error) {
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, err
    }

    fullPath := filepath.Join(logDir, logFileName)
    file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
        currentSize: stat.Size(),
        basePath:    logDir,
    }, nil
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
    if err := lr.currentFile.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
    oldPath := filepath.Join(lr.basePath, logFileName)
    newPath := filepath.Join(lr.basePath, backupName)

    if err := os.Rename(oldPath, newPath); err != nil {
        return err
    }

    file, err := os.OpenFile(oldPath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    lr.currentFile = file
    lr.currentSize = 0

    go lr.cleanupOldFiles()

    return nil
}

func (lr *LogRotator) cleanupOldFiles() {
    pattern := filepath.Join(lr.basePath, logFileName+".*")
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackupFiles {
        return
    }

    filesToDelete := matches[:len(matches)-maxBackupFiles]
    for _, file := range filesToDelete {
        os.Remove(file)
    }
}

func (lr *LogRotator) Close() error {
    return lr.currentFile.Close()
}

func main() {
    rotator, err := NewLogRotator("./logs")
    if err != nil {
        panic(err)
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := rotator.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }
}package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type LogRotator struct {
	CurrentLogPath string
	MaxSize        int64
	ArchiveDir     string
}

func NewLogRotator(logPath string, maxSize int64, archiveDir string) *LogRotator {
	return &LogRotator{
		CurrentLogPath: logPath,
		MaxSize:        maxSize,
		ArchiveDir:     archiveDir,
	}
}

func (lr *LogRotator) CheckAndRotate() error {
	fileInfo, err := os.Stat(lr.CurrentLogPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to stat log file: %w", err)
	}

	if fileInfo.Size() < lr.MaxSize {
		return nil
	}

	timestamp := time.Now().Format("20060102_150405")
	archiveName := fmt.Sprintf("log_%s.log", timestamp)
	archivePath := filepath.Join(lr.ArchiveDir, archiveName)

	err = os.Rename(lr.CurrentLogPath, archivePath)
	if err != nil {
		return fmt.Errorf("failed to rotate log: %w", err)
	}

	newFile, err := os.Create(lr.CurrentLogPath)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}
	newFile.Close()

	fmt.Printf("Log rotated: %s -> %s\n", lr.CurrentLogPath, archivePath)
	return nil
}

func main() {
	rotator := NewLogRotator("app.log", 1024*1024, "./archive")
	err := rotator.CheckAndRotate()
	if err != nil {
		fmt.Printf("Error during log rotation: %v\n", err)
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
    "sync"
    "time"
)

type RotatingLogger struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    currentSize int64
    currentFile *os.File
    fileCounter int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    logger := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSize,
    }
    err := logger.openCurrentFile()
    return logger, err
}

func (rl *RotatingLogger) openCurrentFile() error {
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

func (rl *RotatingLogger) rotate() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentFile == nil {
        return nil
    }

    rl.currentFile.Close()
    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s.%d", rl.basePath, timestamp, rl.fileCounter)
    rl.fileCounter++

    err := os.Rename(rl.basePath, rotatedPath)
    if err != nil {
        return err
    }

    err = rl.compressFile(rotatedPath)
    if err != nil {
        fmt.Printf("Compression failed: %v\n", err)
    }

    return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(sourcePath string) error {
    sourceFile, err := os.Open(sourcePath)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    compressedPath := sourcePath + ".gz"
    compressedFile, err := os.Create(compressedPath)
    if err != nil {
        return err
    }
    defer compressedFile.Close()

    gzWriter := gzip.NewWriter(compressedFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, sourceFile)
    if err != nil {
        return err
    }

    os.Remove(sourcePath)
    return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > rl.maxSize {
        rl.mu.Unlock()
        err := rl.rotate()
        rl.mu.Lock()
        if err != nil {
            return 0, err
        }
    }

    n, err := rl.currentFile.Write(p)
    if err == nil {
        rl.currentSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func (rl *RotatingLogger) cleanupOldFiles(maxFiles int) error {
    dir := filepath.Dir(rl.basePath)
    baseName := filepath.Base(rl.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }

    var matchingFiles []string
    for _, entry := range entries {
        if !entry.IsDir() && strings.HasPrefix(entry.Name(), baseName+".") {
            matchingFiles = append(matchingFiles, entry.Name())
        }
    }

    if len(matchingFiles) <= maxFiles {
        return nil
    }

    sortFilesByAge(matchingFiles)

    for i := 0; i < len(matchingFiles)-maxFiles; i++ {
        filePath := filepath.Join(dir, matchingFiles[i])
        os.Remove(filePath)
    }

    return nil
}

func sortFilesByAge(files []string) {
    extractTimestamp := func(filename string) time.Time {
        parts := strings.Split(filename, ".")
        if len(parts) < 3 {
            return time.Time{}
        }
        timestampStr := parts[len(parts)-2]
        t, err := time.Parse("20060102_150405", timestampStr)
        if err != nil {
            return time.Time{}
        }
        return t
    }

    for i := 0; i < len(files); i++ {
        for j := i + 1; j < len(files); j++ {
            if extractTimestamp(files[i]).After(extractTimestamp(files[j])) {
                files[i], files[j] = files[j], files[i]
            }
        }
    }
}

func main() {
    logger, err := NewRotatingLogger("app.log", 10)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        _, err := logger.Write([]byte(logEntry))
        if err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }

    logger.cleanupOldFiles(5)
    fmt.Println("Log rotation completed")
}