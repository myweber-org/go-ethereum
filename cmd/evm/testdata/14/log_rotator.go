package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	maxFileSize  = 10 * 1024 * 1024 // 10MB
	maxBackupCount = 5
	logFileName   = "app.log"
)

type LogRotator struct {
	currentFile *os.File
	currentSize int64
	basePath    string
}

func NewLogRotator(baseDir string) (*LogRotator, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	fullPath := filepath.Join(baseDir, logFileName)
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to stat log file: %w", err)
	}

	return &LogRotator{
		currentFile: file,
		currentSize: info.Size(),
		basePath:    baseDir,
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
		return fmt.Errorf("failed to close current log file: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
	backupPath := filepath.Join(lr.basePath, backupName)
	originalPath := filepath.Join(lr.basePath, logFileName)

	if err := os.Rename(originalPath, backupPath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	file, err := os.OpenFile(originalPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}

	lr.currentFile = file
	lr.currentSize = 0

	go lr.cleanupOldBackups()

	return nil
}

func (lr *LogRotator) cleanupOldBackups() {
	pattern := filepath.Join(lr.basePath, logFileName+".*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= maxBackupCount {
		return
	}

	sort.Sort(sort.Reverse(sort.StringSlice(matches)))

	for i := maxBackupCount; i < len(matches); i++ {
		os.Remove(matches[i])
	}
}

func (lr *LogRotator) Close() error {
	if lr.currentFile != nil {
		return lr.currentFile.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator("./logs")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating log rotator: %v\n", err)
		os.Exit(1)
	}
	defer rotator.Close()

	messages := []string{
		"Starting application...\n",
		"Processing request from 192.168.1.100\n",
		"Database connection established\n",
		"User login successful\n",
		"Cache updated with new entries\n",
		"Background task completed\n",
		"Sending notification email\n",
		"API response time: 45ms\n",
		"Memory usage: 245MB\n",
		"Shutting down gracefully...\n",
	}

	for _, msg := range messages {
		if _, err := rotator.Write([]byte(msg)); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing log: %v\n", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed. Check ./logs directory")
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

const (
    maxFileSize   = 10 * 1024 * 1024 // 10MB
    maxBackupFiles = 5
)

type LogRotator struct {
    mu          sync.Mutex
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
    lr.mu.Lock()
    defer lr.mu.Unlock()

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
    }

    oldPath := lr.basePath
    if lr.fileIndex > 0 {
        oldPath = fmt.Sprintf("%s.%d", lr.basePath, lr.fileIndex)
    }

    compressedPath := fmt.Sprintf("%s.gz", oldPath)
    if err := compressFile(oldPath, compressedPath); err != nil {
        return err
    }

    os.Remove(oldPath)

    lr.fileIndex++
    if lr.fileIndex > maxBackupFiles {
        lr.fileIndex = 1
    }

    return lr.openCurrentFile()
}

func compressFile(src, dst string) error {
    source, err := os.Open(src)
    if err != nil {
        return err
    }
    defer source.Close()

    destination, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer destination.Close()

    gz := gzip.NewWriter(destination)
    defer gz.Close()

    _, err = io.Copy(gz, source)
    return err
}

func (lr *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(lr.basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

func (lr *LogRotator) Close() error {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func scanExistingLogs(basePath string) int {
    maxIndex := 0
    pattern := fmt.Sprintf("%s.*.gz", filepath.Base(basePath))
    files, _ := filepath.Glob(filepath.Join(filepath.Dir(basePath), pattern))

    for _, file := range files {
        parts := strings.Split(filepath.Base(file), ".")
        if len(parts) >= 3 {
            if idx, err := strconv.Atoi(parts[len(parts)-2]); err == nil && idx > maxIndex {
                maxIndex = idx
            }
        }
    }
    return maxIndex
}

func main() {
    logPath := "./application.log"
    rotator, err := NewLogRotator(logPath)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 100; i++ {
        message := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        rotator.Write([]byte(message))
        time.Sleep(100 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
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
    filename   string
    current    *os.File
    currentSize int64
}

func NewLogRotator(filename string) (*LogRotator, error) {
    rotator := &LogRotator{filename: filename}
    if err := rotator.openCurrent(); err != nil {
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

    n, err := lr.current.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if lr.current != nil {
        lr.current.Close()
    }

    timestamp := time.Now().Format("20060102150405")
    rotatedName := fmt.Sprintf("%s.%s", lr.filename, timestamp)
    if err := os.Rename(lr.filename, rotatedName); err != nil {
        return err
    }

    if err := lr.compressFile(rotatedName); err != nil {
        return err
    }

    if err := lr.cleanupOldBackups(); err != nil {
        return err
    }

    return lr.openCurrent()
}

func (lr *LogRotator) compressFile(filename string) error {
    src, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(filename + ".gz")
    if err != nil {
        return err
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    if _, err := io.Copy(gz, src); err != nil {
        return err
    }

    if err := os.Remove(filename); err != nil {
        return err
    }

    return nil
}

func (lr *LogRotator) cleanupOldBackups() error {
    pattern := lr.filename + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= maxBackups {
        return nil
    }

    var timestamps []time.Time
    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 3 {
            continue
        }
        ts, err := time.Parse("20060102150405", parts[len(parts)-2])
        if err != nil {
            continue
        }
        timestamps = append(timestamps, ts)
    }

    if len(timestamps) > maxBackups {
        oldestIndex := 0
        for i := 1; i < len(timestamps); i++ {
            if timestamps[i].Before(timestamps[oldestIndex]) {
                oldestIndex = i
            }
        }
        if err := os.Remove(matches[oldestIndex]); err != nil {
            return err
        }
    }

    return nil
}

func (lr *LogRotator) openCurrent() error {
    file, err := os.OpenFile(lr.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    lr.current = file
    lr.currentSize = info.Size()
    return nil
}

func (lr *LogRotator) Close() error {
    if lr.current != nil {
        return lr.current.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("application.log")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
    "strings"
    "time"
)

const (
    maxFileSize  = 10 * 1024 * 1024 // 10MB
    maxBackupCount = 5
    logFileName   = "app.log"
)

type LogRotator struct {
    currentFile *os.File
    currentSize int64
    basePath    string
}

func NewLogRotator(baseDir string) (*LogRotator, error) {
    if err := os.MkdirAll(baseDir, 0755); err != nil {
        return nil, err
    }
    
    lr := &LogRotator{
        basePath: baseDir,
    }
    
    if err := lr.openCurrentFile(); err != nil {
        return nil, err
    }
    
    return lr, nil
}

func (lr *LogRotator) openCurrentFile() error {
    fullPath := filepath.Join(lr.basePath, logFileName)
    file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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
    }
    
    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
    oldPath := filepath.Join(lr.basePath, logFileName)
    newPath := filepath.Join(lr.basePath, backupName)
    
    if err := os.Rename(oldPath, newPath); err != nil {
        return err
    }
    
    if err := lr.openCurrentFile(); err != nil {
        return err
    }
    
    lr.cleanupOldBackups()
    return nil
}

func (lr *LogRotator) cleanupOldBackups() {
    files, err := filepath.Glob(filepath.Join(lr.basePath, logFileName+".*"))
    if err != nil {
        return
    }
    
    sort.Sort(sort.Reverse(sort.StringSlice(files)))
    
    for i := maxBackupCount; i < len(files); i++ {
        os.Remove(files[i])
    }
}

func (lr *LogRotator) Close() error {
    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("./logs")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()
    
    for i := 0; i < 100; i++ {
        message := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        
        if i%10 == 0 {
            time.Sleep(100 * time.Millisecond)
        }
    }
    
    fmt.Println("Log rotation test completed")
}