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
	mu            sync.Mutex
	currentFile   *os.File
	currentSize   int64
	maxSize       int64
	logDir        string
	baseName      string
	rotationCount int
	maxRotations  int
}

func NewRotatingLogger(dir, name string, maxSize int64, maxRotations int) (*RotatingLogger, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	rl := &RotatingLogger{
		maxSize:      maxSize,
		logDir:       dir,
		baseName:     name,
		maxRotations: maxRotations,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	path := filepath.Join(rl.logDir, rl.baseName+".log")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to stat log file: %w", err)
	}

	rl.currentFile = file
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) rotate() error {
	rl.currentFile.Close()

	oldPath := filepath.Join(rl.logDir, rl.baseName+".log")
	timestamp := time.Now().Format("20060102_150405")
	rotatedPath := filepath.Join(rl.logDir, fmt.Sprintf("%s_%s.log", rl.baseName, timestamp))

	if err := os.Rename(oldPath, rotatedPath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	compressedPath := rotatedPath + ".gz"
	if err := compressFile(rotatedPath, compressedPath); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to compress log file: %v\n", err)
	} else {
		os.Remove(rotatedPath)
	}

	rl.rotationCount++
	if rl.rotationCount > rl.maxRotations {
		rl.cleanOldRotations()
	}

	return rl.openCurrentFile()
}

func compressFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dest.Close()

	gz := gzip.NewWriter(dest)
	defer gz.Close()

	_, err = io.Copy(gz, source)
	return err
}

func (rl *RotatingLogger) cleanOldRotations() {
	pattern := filepath.Join(rl.logDir, rl.baseName+"_*.log.gz")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) > rl.maxRotations {
		toDelete := matches[:len(matches)-rl.maxRotations]
		for _, file := range toDelete {
			os.Remove(file)
		}
	}
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

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("./logs", "app", 1024*1024, 5)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
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
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type RotatingLogger struct {
    currentFile *os.File
    currentSize int64
    basePath    string
    fileCount   int
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
    rl := &RotatingLogger{basePath: path}
    if err := rl.openCurrentFile(); err != nil {
        return nil, err
    }
    return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    if rl.currentSize+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.currentFile.Write(p)
    rl.currentSize += int64(n)
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.currentFile.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    archivedPath := fmt.Sprintf("%s.%s.gz", rl.basePath, timestamp)
    if err := compressFile(rl.basePath, archivedPath); err != nil {
        return err
    }

    rl.fileCount++
    if rl.fileCount > maxBackups {
        if err := rl.cleanupOldFiles(); err != nil {
            return err
        }
    }

    os.Remove(rl.basePath)
    return rl.openCurrentFile()
}

func compressFile(src, dst string) error {
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()

    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    gz := gzip.NewWriter(out)
    defer gz.Close()

    _, err = io.Copy(gz, in)
    return err
}

func (rl *RotatingLogger) cleanupOldFiles() error {
    pattern := rl.basePath + ".*.gz"
    files, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(files) > maxBackups {
        filesToDelete := files[:len(files)-maxBackups]
        for _, f := range filesToDelete {
            os.Remove(f)
        }
    }
    return nil
}

func (rl *RotatingLogger) openCurrentFile() error {
    f, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := f.Stat()
    if err != nil {
        f.Close()
        return err
    }

    rl.currentFile = f
    rl.currentSize = info.Size()
    return nil
}

func (rl *RotatingLogger) Close() error {
    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}