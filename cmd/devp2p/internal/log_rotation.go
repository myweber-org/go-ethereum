package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

const maxLogSize = 1024 * 1024 // 1MB
const backupCount = 5

type RotatingWriter struct {
	currentSize int64
	basePath    string
	file        *os.File
}

func NewRotatingWriter(path string) (*RotatingWriter, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &RotatingWriter{
		currentSize: info.Size(),
		basePath:    path,
		file:        file,
	}, nil
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
	if w.currentSize+int64(len(p)) > maxLogSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := w.file.Write(p)
	if err == nil {
		w.currentSize += int64(n)
	}
	return n, err
}

func (w *RotatingWriter) rotate() error {
	if err := w.file.Close(); err != nil {
		return err
	}

	for i := backupCount - 1; i >= 0; i-- {
		oldPath := w.backupPath(i)
		newPath := w.backupPath(i + 1)

		if _, err := os.Stat(oldPath); err == nil {
			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
		}
	}

	if err := os.Rename(w.basePath, w.backupPath(0)); err != nil && !os.IsNotExist(err) {
		return err
	}

	file, err := os.OpenFile(w.basePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	w.file = file
	w.currentSize = 0
	return nil
}

func (w *RotatingWriter) backupPath(index int) string {
	if index == 0 {
		return w.basePath
	}
	return w.basePath + "." + string(rune('0'+index))
}

func (w *RotatingWriter) Close() error {
	return w.file.Close()
}

func main() {
	writer, err := NewRotatingWriter("logs/app.log")
	if err != nil {
		log.Fatal(err)
	}
	defer writer.Close()

	log.SetOutput(io.MultiWriter(os.Stdout, writer))

	for i := 0; i < 10000; i++ {
		log.Printf("Log entry number %d", i)
	}
}package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

type RotatingLogger struct {
    currentFile   *os.File
    basePath      string
    maxSize       int64
    rotationCount int
    currentSize   int64
}

func NewRotatingLogger(basePath string, maxSize int64) (*RotatingLogger, error) {
    logger := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSize,
    }
    
    if err := logger.openCurrentFile(); err != nil {
        return nil, err
    }
    
    return logger, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }
    
    filename := fmt.Sprintf("%s.%d.log", rl.basePath, rl.rotationCount)
    file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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
    if rl.currentSize >= rl.maxSize {
        rl.rotationCount++
        return rl.openCurrentFile()
    }
    return nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    if err := rl.rotateIfNeeded(); err != nil {
        return 0, err
    }
    
    n, err = rl.currentFile.Write(p)
    if err == nil {
        rl.currentSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) Close() error {
    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func (rl *RotatingLogger) CleanupOldLogs(maxAge time.Duration) error {
    files, err := filepath.Glob(rl.basePath + ".*.log")
    if err != nil {
        return err
    }
    
    cutoff := time.Now().Add(-maxAge)
    for _, file := range files {
        info, err := os.Stat(file)
        if err != nil {
            continue
        }
        
        if info.ModTime().Before(cutoff) {
            os.Remove(file)
        }
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app", 1024*1024) // 1MB max size
    if err != nil {
        panic(err)
    }
    defer logger.Close()
    
    go func() {
        ticker := time.NewTicker(24 * time.Hour)
        defer ticker.Stop()
        for range ticker.C {
            logger.CleanupOldLogs(7 * 24 * time.Hour) // Keep logs for 7 days
        }
    }()
    
    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(msg))
        time.Sleep(100 * time.Millisecond)
    }
}