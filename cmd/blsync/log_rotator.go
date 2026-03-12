
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
	maxFileSize = 10 * 1024 * 1024
	logDir      = "./logs"
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	mu          sync.Mutex
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
	oldPath := filepath.Join(logDir, fmt.Sprintf("%s_%d.log", rl.baseName, seq))
	newPath := filepath.Join(logDir, fmt.Sprintf("%s_%d.log.gz", rl.baseName, seq))

	oldFile, err := os.Open(oldPath)
	if err != nil {
		return
	}
	defer oldFile.Close()

	newFile, err := os.Create(newPath)
	if err != nil {
		return
	}
	defer newFile.Close()

	gzWriter := gzip.NewWriter(newFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, oldFile); err != nil {
		return
	}

	os.Remove(oldPath)
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
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n",
			time.Now().Format(time.RFC3339), i)
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
	"sync"
	"time"
)

type RotatingLogger struct {
	mu          sync.Mutex
	currentFile *os.File
	basePath    string
	maxSize     int64
	fileSize    int64
	fileCount   int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: basePath,
		maxSize:  int64(maxSizeMB) * 1024 * 1024,
	}
	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	rl.fileCount++
	filename := fmt.Sprintf("%s_%d.log", rl.basePath, rl.fileCount)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	rl.currentFile = file
	info, _ := file.Stat()
	rl.fileSize = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.fileSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.currentFile.Write(p)
	rl.fileSize += int64(n)
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.currentFile.Close(); err != nil {
		return err
	}

	oldPath := fmt.Sprintf("%s_%d.log", rl.basePath, rl.fileCount)
	compressedPath := oldPath + ".gz"

	if err := compressFile(oldPath, compressedPath); err != nil {
		return err
	}

	if err := os.Remove(oldPath); err != nil {
		return err
	}

	return rl.openCurrentFile()
}

func compressFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
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
	return rl.currentFile.Close()
}

func main() {
	logger, err := NewRotatingLogger("app", 10)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n",
			time.Now().Format(time.RFC3339), i)
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
    filename    string
    currentSize int64
    file        *os.File
}

func NewLogRotator(filename string) (*LogRotator, error) {
    rotator := &LogRotator{filename: filename}
    if err := rotator.openFile(); err != nil {
        return nil, err
    }
    return rotator, nil
}

func (lr *LogRotator) openFile() error {
    file, err := os.OpenFile(lr.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    lr.file = file
    lr.currentSize = stat.Size()
    return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }
    n, err := lr.file.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if err := lr.file.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102150405")
    backupName := fmt.Sprintf("%s.%s", lr.filename, timestamp)
    if err := os.Rename(lr.filename, backupName); err != nil {
        return err
    }

    if err := lr.openFile(); err != nil {
        return err
    }

    go lr.compressBackup(backupName)
    go lr.cleanupOldBackups()

    return nil
}

func (lr *LogRotator) compressBackup(filename string) {
    src, err := os.Open(filename)
    if err != nil {
        return
    }
    defer src.Close()

    dst, err := os.Create(filename + ".gz")
    if err != nil {
        return
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    if _, err := io.Copy(gz, src); err != nil {
        return
    }

    os.Remove(filename)
}

func (lr *LogRotator) cleanupOldBackups() {
    pattern := lr.filename + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackups {
        return
    }

    var backupFiles []struct {
        path string
        time time.Time
    }

    for _, match := range matches {
        parts := strings.Split(filepath.Base(match), ".")
        if len(parts) < 3 {
            continue
        }
        timestamp := parts[len(parts)-2]
        t, err := time.Parse("20060102150405", timestamp)
        if err != nil {
            continue
        }
        backupFiles = append(backupFiles, struct {
            path string
            time time.Time
        }{match, t})
    }

    for i := 0; i < len(backupFiles)-maxBackups; i++ {
        oldestIndex := 0
        for j := 1; j < len(backupFiles); j++ {
            if backupFiles[j].time.Before(backupFiles[oldestIndex].time) {
                oldestIndex = j
            }
        }
        os.Remove(backupFiles[oldestIndex].path)
        backupFiles = append(backupFiles[:oldestIndex], backupFiles[oldestIndex+1:]...)
    }
}

func (lr *LogRotator) Close() error {
    if lr.file != nil {
        return lr.file.Close()
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
}