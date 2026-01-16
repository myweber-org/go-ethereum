
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	maxFileSize  = 10 * 1024 * 1024 // 10MB
	maxBackupCount = 5
	logFileName   = "app.log"
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	basePath    string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	fullPath := filepath.Join(path, logFileName)
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
		basePath:    path,
	}, nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	if rl.currentSize+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = rl.currentFile.Write(p)
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
	backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
	backupPath := filepath.Join(rl.basePath, backupName)

	if err := os.Rename(filepath.Join(rl.basePath, logFileName), backupPath); err != nil {
		return err
	}

	file, err := os.OpenFile(filepath.Join(rl.basePath, logFileName), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.currentFile = file
	rl.currentSize = 0

	go rl.cleanupOldBackups()

	return nil
}

func (rl *RotatingLogger) cleanupOldBackups() {
	pattern := filepath.Join(rl.basePath, logFileName+".*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= maxBackupCount {
		return
	}

	for i := 0; i < len(matches)-maxBackupCount; i++ {
		os.Remove(matches[i])
	}
}

func (rl *RotatingLogger) Close() error {
	return rl.currentFile.Close()
}

func main() {
	logger, err := NewRotatingLogger(".")
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(100 * time.Millisecond)
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
    currentFile *os.File
    currentSize int64
    maxFiles    int
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

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    l.currentFile = file
    l.currentSize = info.Size()
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

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", l.basePath, timestamp)

    if err := os.Rename(l.basePath, rotatedPath); err != nil {
        return err
    }

    if err := l.compressFile(rotatedPath); err != nil {
        return err
    }

    if err := l.cleanupOldFiles(); err != nil {
        return err
    }

    return l.openCurrentFile()
}

func (l *RotatingLogger) compressFile(sourcePath string) error {
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

    if _, err := io.Copy(gzWriter, sourceFile); err != nil {
        return err
    }

    os.Remove(sourcePath)
    return nil
}

func (l *RotatingLogger) cleanupOldFiles() error {
    if l.maxFiles <= 0 {
        return nil
    }

    dir := filepath.Dir(l.basePath)
    baseName := filepath.Base(l.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }

    var compressedFiles []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, baseName+".") && strings.HasSuffix(name, ".gz") {
            compressedFiles = append(compressedFiles, name)
        }
    }

    if len(compressedFiles) <= l.maxFiles {
        return nil
    }

    sortByTimestamp := func(files []string) {
        for i := 0; i < len(files); i++ {
            for j := i + 1; j < len(files); j++ {
                timeI := extractTimestamp(files[i])
                timeJ := extractTimestamp(files[j])
                if timeI.After(timeJ) {
                    files[i], files[j] = files[j], files[i]
                }
            }
        }
    }

    sortByTimestamp(compressedFiles)

    filesToRemove := compressedFiles[:len(compressedFiles)-l.maxFiles]
    for _, file := range filesToRemove {
        os.Remove(filepath.Join(dir, file))
    }

    return nil
}

func extractTimestamp(filename string) time.Time {
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
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n",
            time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}