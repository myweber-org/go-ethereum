package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxFileSize  = 10 * 1024 * 1024
	maxBackups   = 5
	logDirectory = "./logs"
)

type RotatingLogger struct {
	mu        sync.Mutex
	file      *os.File
	size      int64
	basePath  string
	sequence  int
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDirectory, 0755); err != nil {
		return nil, err
	}
	basePath := filepath.Join(logDirectory, baseName)
	rl := &RotatingLogger{basePath: basePath}
	if err := rl.openFile(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openFile() error {
	path := rl.basePath
	if rl.sequence > 0 {
		path = fmt.Sprintf("%s.%d", rl.basePath, rl.sequence)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}
	rl.file = file
	rl.size = info.Size()
	return nil
}

func (rl *RotatingLogger) rotate() error {
	rl.file.Close()
	rl.sequence++
	if rl.sequence > maxBackups {
		rl.sequence = maxBackups
		oldest := fmt.Sprintf("%s.%d", rl.basePath, rl.sequence)
		if err := os.Remove(oldest); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	for i := rl.sequence; i > 1; i-- {
		oldName := fmt.Sprintf("%s.%d", rl.basePath, i-1)
		newName := fmt.Sprintf("%s.%d", rl.basePath, i)
		if err := os.Rename(oldName, newName); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	if err := os.Rename(rl.basePath, rl.basePath+".1"); err != nil {
		return err
	}
	rl.sequence = 0
	return rl.openFile()
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if rl.size+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}
	n, err = rl.file.Write(p)
	rl.size += int64(n)
	return n, err
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func main() {
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, logger))
	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d at %s", i, time.Now().Format(time.RFC3339))
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
    file        *os.File
    currentSize int64
    maxSize     int64
    basePath    string
    fileCount   int
    maxFiles    int
}

func NewRotatingLogger(basePath string, maxSize int64, maxFiles int) (*RotatingLogger, error) {
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
        currentSize: info.Size(),
        maxSize:     maxSize,
        basePath:    basePath,
        maxFiles:    maxFiles,
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

    timestamp := time.Now().Format("20060102_150405")
    archivePath := fmt.Sprintf("%s.%s.gz", rl.basePath, timestamp)

    if err := compressFile(rl.basePath, archivePath); err != nil {
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
    rl.fileCount++

    if rl.fileCount > rl.maxFiles {
        rl.cleanupOldFiles()
    }

    return nil
}

func compressFile(source, target string) error {
    srcFile, err := os.Open(source)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    destFile, err := os.Create(target)
    if err != nil {
        return err
    }
    defer destFile.Close()

    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, srcFile)
    return err
}

func (rl *RotatingLogger) cleanupOldFiles() {
    pattern := rl.basePath + ".*.gz"
    files, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(files) > rl.maxFiles {
        filesToDelete := files[:len(files)-rl.maxFiles]
        for _, file := range filesToDelete {
            os.Remove(file)
        }
    }
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.file.Close()
}

func main() {
    logger, err := NewRotatingLogger("app.log", 1024*1024, 5)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Some sample data here\n",
            time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }
}