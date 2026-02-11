package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
)

const maxFileSize = 1024 * 1024 // 1MB

func rotateLogFile(filePath string) error {
    fileInfo, err := os.Stat(filePath)
    if err != nil {
        return err
    }

    if fileInfo.Size() < maxFileSize {
        return nil
    }

    for i := 1; ; i++ {
        newPath := filePath + "." + strconv.Itoa(i)
        if _, err := os.Stat(newPath); os.IsNotExist(err) {
            err := os.Rename(filePath, newPath)
            if err != nil {
                return err
            }
            fmt.Printf("Rotated log file to: %s\n", newPath)
            break
        }
    }

    newFile, err := os.Create(filePath)
    if err != nil {
        return err
    }
    newFile.Close()
    return nil
}

func writeLog(filePath, message string) error {
    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = io.WriteString(file, message+"\n")
    if err != nil {
        return err
    }

    return rotateLogFile(filePath)
}

func main() {
    logPath := filepath.Join(".", "application.log")
    for i := 0; i < 15000; i++ {
        err := writeLog(logPath, fmt.Sprintf("Log entry number: %d", i))
        if err != nil {
            fmt.Printf("Error writing log: %v\n", err)
            break
        }
    }
}package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotatingLogger struct {
	mu          sync.Mutex
	file        *os.File
	basePath    string
	maxSize     int64
	currentSize int64
	rotationNum int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

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
		file:        file,
		basePath:    basePath,
		maxSize:     maxSize,
		currentSize: info.Size(),
		rotationNum: 0,
	}, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.file.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.file.Close(); err != nil {
		return err
	}

	rl.rotationNum++
	archiveName := fmt.Sprintf("%s.%d-%s.gz",
		rl.basePath,
		rl.rotationNum,
		time.Now().Format("20060102-150405"))

	if err := rl.compressFile(rl.basePath, archiveName); err != nil {
		return err
	}

	if err := os.Remove(rl.basePath); err != nil {
		return err
	}

	file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.file = file
	rl.currentSize = 0
	return nil
}

func (rl *RotatingLogger) compressFile(source, target string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	targetDir := filepath.Dir(target)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

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

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func main() {
	logger, err := NewRotatingLogger("app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}