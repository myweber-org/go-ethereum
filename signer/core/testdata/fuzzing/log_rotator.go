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
    maxBackups   = 5
    logDirectory = "./logs"
)

type LogRotator struct {
    currentFile *os.File
    filePath    string
    bytesWritten int64
}

func NewLogRotator(baseName string) (*LogRotator, error) {
    if err := os.MkdirAll(logDirectory, 0755); err != nil {
        return nil, err
    }
    
    filePath := filepath.Join(logDirectory, baseName+".log")
    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    
    info, _ := file.Stat()
    return &LogRotator{
        currentFile: file,
        filePath:    filePath,
        bytesWritten: info.Size(),
    }, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.bytesWritten+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }
    
    n, err := lr.currentFile.Write(p)
    if err == nil {
        lr.bytesWritten += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    lr.currentFile.Close()
    
    timestamp := time.Now().Format("20060102-150405")
    rotatedPath := filepath.Join(logDirectory, fmt.Sprintf("%s.%s.log", 
        filepath.Base(lr.filePath[:len(lr.filePath)-4]), timestamp))
    
    if err := os.Rename(lr.filePath, rotatedPath); err != nil {
        return err
    }
    
    file, err := os.OpenFile(lr.filePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    
    lr.currentFile = file
    lr.bytesWritten = 0
    
    go lr.cleanupOldFiles()
    return nil
}

func (lr *LogRotator) cleanupOldFiles() {
    pattern := filepath.Join(logDirectory, filepath.Base(lr.filePath[:len(lr.filePath)-4])+".*.log")
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
    rotator, err := NewLogRotator("application")
    if err != nil {
        panic(err)
    }
    defer rotator.Close()
    
    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n", 
            time.Now().Format(time.RFC3339), i)
        rotator.Write([]byte(message))
        time.Sleep(10 * time.Millisecond)
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
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	logDir      = "./logs"
)

type RotatingLogger struct {
	mu          sync.Mutex
	currentFile *os.File
	currentSize int64
	baseName    string
	sequence    int
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: baseName,
		sequence: 0,
	}

	if err := rl.openNewFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openNewFile() error {
	rl.sequence++
	filename := filepath.Join(logDir, fmt.Sprintf("%s_%d.log", rl.baseName, rl.sequence))
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if rl.currentFile != nil {
		rl.currentFile.Close()
		go rl.compressPreviousFile(rl.sequence - 1)
	}

	rl.currentFile = file
	rl.currentSize = 0
	return nil
}

func (rl *RotatingLogger) compressPreviousFile(seq int) {
	oldFile := filepath.Join(logDir, fmt.Sprintf("%s_%d.log", rl.baseName, seq))
	compressedFile := filepath.Join(logDir, fmt.Sprintf("%s_%d.log.gz", rl.baseName, seq))

	src, err := os.Open(oldFile)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.Create(compressedFile)
	if err != nil {
		return
	}
	defer dst.Close()

	gz := gzip.NewWriter(dst)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err != nil {
		return
	}

	os.Remove(oldFile)
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > maxFileSize {
		if err := rl.openNewFile(); err != nil {
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

func main() {
	logger, err := NewRotatingLogger("app")
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Application event recorded\n",
			time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(msg))
		time.Sleep(100 * time.Millisecond)
	}
}