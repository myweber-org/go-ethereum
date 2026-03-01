
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type RotatingLogger struct {
	mu           sync.Mutex
	basePath     string
	currentFile  *os.File
	maxSize      int64
	currentSize  int64
	backupCount  int
	compressOld  bool
}

func NewRotatingLogger(basePath string, maxSizeMB int, backupCount int, compressOld bool) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	rl := &RotatingLogger{
		basePath:    basePath,
		maxSize:     maxSize,
		backupCount: backupCount,
		compressOld: compressOld,
	}

	if err := rl.rotateIfNeeded(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile == nil {
		if err := rl.openCurrentFile(); err != nil {
			return 0, err
		}
	}

	if rl.currentSize+int64(len(p)) > rl.maxSize {
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

func (rl *RotatingLogger) openCurrentFile() error {
	dir := filepath.Dir(rl.basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(rl.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

func (rl *RotatingLogger) rotateIfNeeded() error {
	info, err := os.Stat(rl.basePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	if info.Size() >= rl.maxSize {
		return rl.rotate()
	}
	return nil
}

func (rl *RotatingLogger) rotate() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
		rl.currentFile = nil
	}

	for i := rl.backupCount - 1; i >= 0; i-- {
		srcPath := rl.getBackupPath(i)
		dstPath := rl.getBackupPath(i + 1)

		if _, err := os.Stat(srcPath); err == nil {
			if i == rl.backupCount-1 {
				os.Remove(srcPath)
			} else {
				os.Rename(srcPath, dstPath)
			}
		}
	}

	if _, err := os.Stat(rl.basePath); err == nil {
		backupPath := rl.getBackupPath(0)
		os.Rename(rl.basePath, backupPath)

		if rl.compressOld {
			go rl.compressFile(backupPath)
		}
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) getBackupPath(index int) string {
	if index == 0 {
		return rl.basePath + ".1"
	}
	return fmt.Sprintf("%s.%d", rl.basePath, index+1)
}

func (rl *RotatingLogger) compressFile(path string) {
	compressedPath := path + ".gz"
	
	src, err := os.Open(path)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.Create(compressedPath)
	if err != nil {
		return
	}
	defer dst.Close()

	// In production, use compress/gzip here
	// For simplicity, we just copy for this example
	io.Copy(dst, src)
	
	os.Remove(path)
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
	logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10, 5, true)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	customLog := log.New(logger, "", log.LstdFlags)

	for i := 0; i < 1000; i++ {
		customLog.Printf("Log entry %d: %s", i, strings.Repeat("X", 1024))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation completed")
}package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
)

const maxFileSize = 1024 * 1024 // 1MB

type LogRotator struct {
    filePath string
    file     *os.File
    size     int64
}

func NewLogRotator(path string) (*LogRotator, error) {
    file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &LogRotator{
        filePath: path,
        file:     file,
        size:     info.Size(),
    }, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.size+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := lr.file.Write(p)
    if err == nil {
        lr.size += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    lr.file.Close()

    for i := 9; i >= 1; i-- {
        oldName := fmt.Sprintf("%s.%d", lr.filePath, i)
        newName := fmt.Sprintf("%s.%d", lr.filePath, i+1)

        if _, err := os.Stat(oldName); err == nil {
            os.Rename(oldName, newName)
        }
    }

    if _, err := os.Stat(lr.filePath); err == nil {
        backupName := fmt.Sprintf("%s.1", lr.filePath)
        os.Rename(lr.filePath, backupName)
    }

    file, err := os.OpenFile(lr.filePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    lr.file = file
    lr.size = 0
    return nil
}

func (lr *LogRotator) Close() error {
    return lr.file.Close()
}

func main() {
    logPath := filepath.Join(os.TempDir(), "application.log")
    rotator, err := NewLogRotator(logPath)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d: This is a test log message.\n", i)
        if _, err := rotator.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
    }

    fmt.Println("Log rotation test completed")
}