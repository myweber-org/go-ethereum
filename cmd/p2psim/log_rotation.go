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

func NewLogRotator(basePath string) (*LogRotator, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, err
	}

	fullPath := filepath.Join(basePath, logFileName)
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
		basePath:    basePath,
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

	go lr.cleanupOldLogs()
	return nil
}

func (lr *LogRotator) cleanupOldLogs() {
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
	return lr.currentFile.Close()
}

func main() {
	rotator, err := NewLogRotator("./logs")
	if err != nil {
		panic(err)
	}
	defer rotator.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("[%s] Log entry number %d\n", 
			time.Now().Format(time.RFC3339), i)
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
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
    logger := &RotatingLogger{filePath: path}
    if err := logger.openCurrentFile(); err != nil {
        return nil, err
    }
    return logger, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    if rl.bytesWritten+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }
    
    n, err := rl.currentFile.Write(p)
    if err == nil {
        rl.bytesWritten += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }
    
    timestamp := time.Now().Format("20060102_150405")
    archivedPath := fmt.Sprintf("%s.%s.gz", rl.filePath, timestamp)
    
    if err := compressFile(rl.filePath, archivedPath); err != nil {
        return err
    }
    
    if err := cleanupOldBackups(rl.filePath); err != nil {
        return err
    }
    
    os.Remove(rl.filePath)
    return rl.openCurrentFile()
}

func (rl *RotatingLogger) openCurrentFile() error {
    file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    
    rl.currentFile = file
    rl.bytesWritten = info.Size()
    return nil
}

func compressFile(source, target string) error {
    srcFile, err := os.Open(source)
    if err != nil {
        return err
    }
    defer srcFile.Close()
    
    dstFile, err := os.Create(target)
    if err != nil {
        return err
    }
    defer dstFile.Close()
    
    gzWriter := gzip.NewWriter(dstFile)
    defer gzWriter.Close()
    
    _, err = io.Copy(gzWriter, srcFile)
    return err
}

func cleanupOldBackups(basePath string) error {
    pattern := filepath.Join(filepath.Dir(basePath), filepath.Base(basePath)+".*.gz")
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }
    
    if len(matches) > maxBackups {
        filesToRemove := matches[:len(matches)-maxBackups]
        for _, file := range filesToRemove {
            os.Remove(file)
        }
    }
    return nil
}

func (rl *RotatingLogger) Close() error {
    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app.log")
    if err != nil {
        panic(err)
    }
    defer logger.Close()
    
    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Iteration %d: Processing data batch\n", 
            time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }
    
    fmt.Println("Log rotation demonstration completed")
}