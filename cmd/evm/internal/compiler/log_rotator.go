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
	mu           sync.Mutex
	currentFile  *os.File
	basePath     string
	maxSize      int64
	currentSize  int64
	rotationCount int
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
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	file, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	rl.rotationCount++
	archivePath := fmt.Sprintf("%s.%d.%s.gz", 
		rl.basePath, 
		rl.rotationCount, 
		time.Now().Format("20060102_150405"))

	if err := rl.compressFile(rl.basePath, archivePath); err != nil {
		return err
	}

	if err := os.Remove(rl.basePath); err != nil && !os.IsNotExist(err) {
		return err
	}

	return rl.openCurrentFile()
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

	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: Application event occurred\n", 
			time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}
}package main

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
)

type Rotator struct {
    FilePath    string
    MaxSize     int64
    MaxFiles    int
    RotateEvery time.Duration
    lastRotate  time.Time
}

func NewRotator(path string, maxSize int64, maxFiles int, rotateEvery time.Duration) *Rotator {
    return &Rotator{
        FilePath:    path,
        MaxSize:     maxSize,
        MaxFiles:    maxFiles,
        RotateEvery: rotateEvery,
        lastRotate:  time.Now(),
    }
}

func (r *Rotator) Write(p []byte) (int, error) {
    if err := r.maybeRotate(); err != nil {
        return 0, err
    }

    f, err := os.OpenFile(r.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return 0, err
    }
    defer f.Close()

    return f.Write(p)
}

func (r *Rotator) maybeRotate() error {
    now := time.Now()
    shouldRotate := false

    if info, err := os.Stat(r.FilePath); err == nil {
        if info.Size() >= r.MaxSize {
            shouldRotate = true
        }
    }

    if now.Sub(r.lastRotate) >= r.RotateEvery {
        shouldRotate = true
    }

    if shouldRotate {
        if err := r.performRotation(); err != nil {
            return err
        }
        r.lastRotate = now
    }
    return nil
}

func (r *Rotator) performRotation() error {
    for i := r.MaxFiles - 1; i > 0; i-- {
        oldName := fmt.Sprintf("%s.%d", r.FilePath, i)
        newName := fmt.Sprintf("%s.%d", r.FilePath, i+1)

        if _, err := os.Stat(oldName); err == nil {
            os.Rename(oldName, newName)
        }
    }

    if _, err := os.Stat(r.FilePath); err == nil {
        backupName := fmt.Sprintf("%s.1", r.FilePath)
        os.Rename(r.FilePath, backupName)
    }

    return nil
}

func main() {
    rotator := NewRotator(
        filepath.Join("logs", "app.log"),
        10*1024*1024,
        5,
        24*time.Hour,
    )

    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        rotator.Write([]byte(msg))
        time.Sleep(100 * time.Millisecond)
    }

    fmt.Println("Log rotation example completed")
}