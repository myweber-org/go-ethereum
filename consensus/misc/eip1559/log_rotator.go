package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
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

	go lr.cleanupOldLogs()
	return nil
}

func (lr *LogRotator) cleanupOldLogs() {
	pattern := filepath.Join(logDir, lr.baseName+".*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) > maxBackups {
		toDelete := matches[:len(matches)-maxBackups]
		for _, path := range toDelete {
			os.Remove(path)
		}
	}
}

func (lr *LogRotator) Close() error {
	return lr.currentFile.Close()
}

func main() {
	rotator, err := NewLogRotator("app.log")
	if err != nil {
		panic(err)
	}
	defer rotator.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

type Rotator struct {
    filePath    string
    maxSize     int64
    maxAge      time.Duration
    currentFile *os.File
    currentSize int64
}

func NewRotator(filePath string, maxSize int64, maxAge time.Duration) (*Rotator, error) {
    r := &Rotator{
        filePath: filePath,
        maxSize:  maxSize,
        maxAge:   maxAge,
    }
    if err := r.openCurrent(); err != nil {
        return nil, err
    }
    return r, nil
}

func (r *Rotator) openCurrent() error {
    if err := os.MkdirAll(filepath.Dir(r.filePath), 0755); err != nil {
        return err
    }
    f, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    info, err := f.Stat()
    if err != nil {
        f.Close()
        return err
    }
    r.currentFile = f
    r.currentSize = info.Size()
    return nil
}

func (r *Rotator) Write(p []byte) (int, error) {
    if r.needsRotation() {
        if err := r.rotate(); err != nil {
            return 0, err
        }
    }
    n, err := r.currentFile.Write(p)
    r.currentSize += int64(n)
    return n, err
}

func (r *Rotator) needsRotation() bool {
    if r.currentSize >= r.maxSize {
        return true
    }
    if r.maxAge > 0 {
        info, err := os.Stat(r.filePath)
        if err == nil && time.Since(info.ModTime()) > r.maxAge {
            return true
        }
    }
    return false
}

func (r *Rotator) rotate() error {
    if r.currentFile != nil {
        r.currentFile.Close()
    }
    timestamp := time.Now().Format("20060102_150405")
    backupPath := r.filePath + "." + timestamp
    if err := os.Rename(r.filePath, backupPath); err != nil {
        return err
    }
    if err := r.openCurrent(); err != nil {
        return err
    }
    r.cleanOldBackups()
    return nil
}

func (r *Rotator) cleanOldBackups() {
    if r.maxAge <= 0 {
        return
    }
    cutoff := time.Now().Add(-r.maxAge)
    dir := filepath.Dir(r.filePath)
    base := filepath.Base(r.filePath)
    entries, err := os.ReadDir(dir)
    if err != nil {
        return
    }
    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }
        name := entry.Name()
        if len(name) > len(base) && name[:len(base)] == base && name[len(base)] == '.' {
            info, err := entry.Info()
            if err != nil {
                continue
            }
            if info.ModTime().Before(cutoff) {
                os.Remove(filepath.Join(dir, name))
            }
        }
    }
}

func (r *Rotator) Close() error {
    if r.currentFile != nil {
        return r.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewRotator("logs/app.log", 1024*1024, 24*time.Hour)
    if err != nil {
        fmt.Printf("Failed to create rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := rotator.Write([]byte(msg)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        time.Sleep(100 * time.Millisecond)
    }
    fmt.Println("Log rotation example completed")
}