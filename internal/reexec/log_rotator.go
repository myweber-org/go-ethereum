package main

import (
	"fmt"
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
	currentFile *os.File
	currentSize int64
	baseName    string
	mu          sync.Mutex
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: baseName,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(logDir, fmt.Sprintf("%s_%s.log", rl.baseName, timestamp))

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.currentFile = file
	info, _ := file.Stat()
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > maxFileSize {
		rl.currentFile.Close()
		if err := rl.openCurrentFile(); err != nil {
			return 0, err
		}
		rl.currentSize = 0
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
	return rl.currentFile.Close()
}

func main() {
	logger, err := NewRotatingLogger("app")
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		logger.Write([]byte(message))
		time.Sleep(100 * time.Millisecond)
	}
}
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	maxFileSize = 1024 * 1024 * 10 // 10MB
	maxBackups  = 5
	logDir      = "./logs"
)

type RotatingWriter struct {
	currentFile *os.File
	currentSize int64
	baseName    string
}

func NewRotatingWriter(name string) (*RotatingWriter, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	basePath := filepath.Join(logDir, name)
	w := &RotatingWriter{baseName: basePath}

	if err := w.rotateIfNeeded(); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *RotatingWriter) Write(p []byte) (n int, err error) {
	if err := w.rotateIfNeeded(); err != nil {
		return 0, err
	}

	n, err = w.currentFile.Write(p)
	w.currentSize += int64(n)
	return n, err
}

func (w *RotatingWriter) rotateIfNeeded() error {
	if w.currentFile == nil || w.currentSize >= maxFileSize {
		return w.rotate()
	}
	return nil
}

func (w *RotatingWriter) rotate() error {
	if w.currentFile != nil {
		w.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	newPath := fmt.Sprintf("%s_%s.log", w.baseName, timestamp)

	file, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	w.currentFile = file
	w.currentSize = 0

	go w.cleanupOldFiles()

	return nil
}

func (w *RotatingWriter) cleanupOldFiles() {
	pattern := fmt.Sprintf("%s_*.log", w.baseName)
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

func (w *RotatingWriter) Close() error {
	if w.currentFile != nil {
		return w.currentFile.Close()
	}
	return nil
}

func main() {
	writer, err := NewRotatingWriter("app")
	if err != nil {
		panic(err)
	}
	defer writer.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		writer.Write([]byte(msg))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation completed")
}
package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type Rotator struct {
    mu            sync.Mutex
    file          *os.File
    basePath      string
    maxSize       int64
    rotationCount int
    currentSize   int64
    lastRotation  time.Time
    rotationTime  time.Duration
}

func NewRotator(basePath string, maxSize int64, rotationTime time.Duration) (*Rotator, error) {
    dir := filepath.Dir(basePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, fmt.Errorf("create directory: %w", err)
    }

    file, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, fmt.Errorf("open log file: %w", err)
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, fmt.Errorf("stat log file: %w", err)
    }

    return &Rotator{
        file:         file,
        basePath:     basePath,
        maxSize:      maxSize,
        currentSize:  info.Size(),
        lastRotation: time.Now(),
        rotationTime: rotationTime,
    }, nil
}

func (r *Rotator) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if err := r.checkRotation(); err != nil {
        return 0, err
    }

    n, err := r.file.Write(p)
    if err != nil {
        return n, err
    }

    r.currentSize += int64(n)
    return n, nil
}

func (r *Rotator) checkRotation() error {
    now := time.Now()
    shouldRotate := false

    if r.maxSize > 0 && r.currentSize >= r.maxSize {
        shouldRotate = true
    }

    if r.rotationTime > 0 && now.Sub(r.lastRotation) >= r.rotationTime {
        shouldRotate = true
    }

    if !shouldRotate {
        return nil
    }

    if err := r.rotate(); err != nil {
        return fmt.Errorf("rotate log file: %w", err)
    }

    r.lastRotation = now
    return nil
}

func (r *Rotator) rotate() error {
    if err := r.file.Close(); err != nil {
        return fmt.Errorf("close current log file: %w", err)
    }

    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s.%d", r.basePath, timestamp, r.rotationCount)
    r.rotationCount++

    if err := os.Rename(r.basePath, backupPath); err != nil {
        return fmt.Errorf("rename log file: %w", err)
    }

    file, err := os.OpenFile(r.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("create new log file: %w", err)
    }

    r.file = file
    r.currentSize = 0

    return nil
}

func (r *Rotator) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.file != nil {
        return r.file.Close()
    }
    return nil
}

func main() {
    rotator, err := NewRotator("/var/log/myapp/app.log", 10*1024*1024, 24*time.Hour)
    if err != nil {
        fmt.Printf("Failed to create rotator: %v\n", err)
        os.Exit(1)
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