package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	maxFileSize = 1024 * 1024 // 1MB
	maxBackups  = 5
	logDir      = "./logs"
)

type LogRotator struct {
	currentFile *os.File
	currentSize int64
	baseName    string
}

func NewLogRotator(filename string) (*LogRotator, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	basePath := filepath.Join(logDir, filename)
	file, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &LogRotator{
		currentFile: file,
		currentSize: info.Size(),
		baseName:    filename,
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
	oldPath := filepath.Join(logDir, lr.baseName)
	newPath := filepath.Join(logDir, fmt.Sprintf("%s.%s", lr.baseName, timestamp))

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
	pattern := filepath.Join(logDir, lr.baseName+".*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) > maxBackups {
		filesToDelete := matches[:len(matches)-maxBackups]
		for _, file := range filesToDelete {
			os.Remove(file)
		}
	}
}

func (lr *LogRotator) Close() error {
	return lr.currentFile.Close()
}

func main() {
	rotator, err := NewLogRotator("app.log")
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		return
	}
	defer rotator.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}package main

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
    maxFiles    int
    currentFile *os.File
    currentSize int64
}

func NewRotatingLogger(basePath string, maxSizeMB int, maxFiles int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    logger := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSize,
        maxFiles: maxFiles,
    }

    if err := logger.openCurrentFile(); err != nil {
        return nil, err
    }
    return logger, nil
}

func (l *RotatingLogger) openCurrentFile() error {
    file, err := os.OpenFile(l.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    l.currentFile = file
    l.currentSize = stat.Size()
    return nil
}

func (l *RotatingLogger) Write(p []byte) (int, error) {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.currentSize+int64(len(p)) > l.maxSize {
        if err := l.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := l.currentFile.Write(p)
    if err == nil {
        l.currentSize += int64(n)
    }
    return n, err
}

func (l *RotatingLogger) rotate() error {
    if l.currentFile != nil {
        l.currentFile.Close()
    }

    for i := l.maxFiles - 1; i > 0; i-- {
        oldPath := l.getArchivePath(i - 1)
        newPath := l.getArchivePath(i)

        if _, err := os.Stat(oldPath); err == nil {
            if i == l.maxFiles-1 {
                os.Remove(newPath)
            } else {
                os.Rename(oldPath, newPath)
            }
        }
    }

    if err := l.compressCurrentLog(); err != nil {
        return err
    }

    return l.openCurrentFile()
}

func (l *RotatingLogger) compressCurrentLog() error {
    source, err := os.Open(l.basePath)
    if err != nil {
        return err
    }
    defer source.Close()

    timestamp := time.Now().Format("20060102_150405")
    compressedPath := fmt.Sprintf("%s.%s.gz", l.basePath, timestamp)

    dest, err := os.Create(compressedPath)
    if err != nil {
        return err
    }
    defer dest.Close()

    gzWriter := gzip.NewWriter(dest)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, source)
    if err != nil {
        return err
    }

    return os.Remove(l.basePath)
}

func (l *RotatingLogger) getArchivePath(index int) string {
    if index == 0 {
        return l.basePath
    }
    return fmt.Sprintf("%s.%d.gz", l.basePath, index)
}

func (l *RotatingLogger) Close() error {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.currentFile != nil {
        return l.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app.log", 10, 5)
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(message))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}