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
}package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize  = 10 * 1024 * 1024 // 10MB
    maxBackupCount = 5
    logFileName   = "app.log"
)

type LogRotator struct {
    currentSize int64
    basePath    string
}

func NewLogRotator(basePath string) *LogRotator {
    return &LogRotator{
        basePath: basePath,
    }
}

func (lr *LogRotator) Write(p []byte) (n int, err error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    file, err := os.OpenFile(filepath.Join(lr.basePath, logFileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    n, err = file.Write(p)
    lr.currentSize += int64(n)
    return n, err
}

func (lr *LogRotator) rotate() error {
    currentPath := filepath.Join(lr.basePath, logFileName)
    if _, err := os.Stat(currentPath); os.IsNotExist(err) {
        return nil
    }

    timestamp := time.Now().Format("20060102_150405")
    backupPath := filepath.Join(lr.basePath, fmt.Sprintf("%s.%s", logFileName, timestamp))

    if err := os.Rename(currentPath, backupPath); err != nil {
        return err
    }

    lr.currentSize = 0
    lr.cleanupOldBackups()
    return nil
}

func (lr *LogRotator) cleanupOldBackups() {
    files, err := filepath.Glob(filepath.Join(lr.basePath, logFileName+".*"))
    if err != nil {
        return
    }

    sort.Sort(sort.Reverse(sort.StringSlice(files)))

    for i, file := range files {
        if i >= maxBackupCount {
            os.Remove(file)
        }
    }
}

func (lr *LogRotator) initialize() error {
    filePath := filepath.Join(lr.basePath, logFileName)
    info, err := os.Stat(filePath)
    if os.IsNotExist(err) {
        return nil
    }
    if err != nil {
        return err
    }

    lr.currentSize = info.Size()
    return nil
}

func main() {
    rotator := NewLogRotator(".")
    if err := rotator.initialize(); err != nil {
        fmt.Printf("Initialization error: %v\n", err)
        return
    }

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogRotator struct {
	mu          sync.Mutex
	currentFile *os.File
	filePath    string
	maxSize     int64
	currentSize int64
	rotationSeq int
}

func NewLogRotator(basePath string, maxSizeMB int) (*LogRotator, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	rotator := &LogRotator{
		filePath: basePath,
		maxSize:  maxSize,
	}
	if err := rotator.openCurrentFile(); err != nil {
		return nil, err
	}
	return rotator, nil
}

func (lr *LogRotator) openCurrentFile() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentFile != nil {
		lr.currentFile.Close()
	}

	file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	lr.currentFile = file
	lr.currentSize = info.Size()
	return nil
}

func (lr *LogRotator) rotate() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentFile == nil {
		return fmt.Errorf("no current log file")
	}

	lr.currentFile.Close()
	timestamp := time.Now().Format("20060102_150405")
	rotatedPath := fmt.Sprintf("%s.%s.%d", lr.filePath, timestamp, lr.rotationSeq)
	lr.rotationSeq++

	if err := os.Rename(lr.filePath, rotatedPath); err != nil {
		return err
	}

	return lr.openCurrentFile()
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentFile == nil {
		return 0, fmt.Errorf("log file not open")
	}

	if lr.currentSize+int64(len(p)) > lr.maxSize {
		lr.mu.Unlock()
		if err := lr.rotate(); err != nil {
			lr.mu.Lock()
			return 0, err
		}
		lr.mu.Lock()
	}

	n, err := lr.currentFile.Write(p)
	if err == nil {
		lr.currentSize += int64(n)
	}
	return n, err
}

func (lr *LogRotator) Close() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentFile != nil {
		return lr.currentFile.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator("app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		return
	}
	defer rotator.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}