
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
    maxBackups  = 5
)

type RotatingLogger struct {
    filename   string
    current    *os.File
    size       int64
    mu         sync.Mutex
}

func NewRotatingLogger(filename string) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        filename: filename,
    }

    if err := rl.openCurrent(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLogger) openCurrent() error {
    file, err := os.OpenFile(rl.filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    rl.current = file
    rl.size = info.Size()
    return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.size+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.current.Write(p)
    if err == nil {
        rl.size += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.current.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s", rl.filename, timestamp)

    if err := os.Rename(rl.filename, backupName); err != nil {
        return err
    }

    if err := rl.compressBackup(backupName); err != nil {
        return err
    }

    if err := rl.cleanOldBackups(); err != nil {
        return err
    }

    return rl.openCurrent()
}

func (rl *RotatingLogger) compressBackup(filename string) error {
    src, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer src.Close()

    dest, err := os.Create(filename + ".gz")
    if err != nil {
        return err
    }
    defer dest.Close()

    gz := gzip.NewWriter(dest)
    defer gz.Close()

    if _, err := io.Copy(gz, src); err != nil {
        return err
    }

    if err := os.Remove(filename); err != nil {
        return err
    }

    return nil
}

func (rl *RotatingLogger) cleanOldBackups() error {
    pattern := rl.filename + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= maxBackups {
        return nil
    }

    toDelete := matches[:len(matches)-maxBackups]
    for _, file := range toDelete {
        if err := os.Remove(file); err != nil {
            return err
        }
    }

    return nil
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.current != nil {
        return rl.current.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app.log")
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(message))
        time.Sleep(10 * time.Millisecond)
    }
}package main

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

const (
	maxFileSize   = 10 * 1024 * 1024 // 10MB
	backupCount   = 5
	checkInterval = 30 * time.Second
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	filePath   string
	currentPos int64
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{filePath: path}
	if err := rl.openFile(); err != nil {
		return nil, err
	}
	go rl.monitor()
	return rl, nil
}

func (rl *RotatingLogger) openFile() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		rl.file.Close()
	}

	file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.file = file
	rl.currentPos = stat.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	n, err = rl.file.Write(p)
	if err == nil {
		rl.currentPos += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentPos < maxFileSize {
		return nil
	}

	rl.file.Close()

	baseDir := filepath.Dir(rl.filePath)
	baseName := filepath.Base(rl.filePath)
	ext := filepath.Ext(baseName)
	nameWithoutExt := strings.TrimSuffix(baseName, ext)

	for i := backupCount - 1; i >= 0; i-- {
		var src, dst string
		if i == 0 {
			src = rl.filePath
		} else {
			src = filepath.Join(baseDir, fmt.Sprintf("%s.%d%s", nameWithoutExt, i, ext))
		}
		dst = filepath.Join(baseDir, fmt.Sprintf("%s.%d%s", nameWithoutExt, i+1, ext))

		if _, err := os.Stat(src); err == nil {
			if i == backupCount-1 {
				os.Remove(dst)
			} else {
				os.Rename(src, dst)
			}
		}
	}

	return rl.openFile()
}

func (rl *RotatingLogger) monitor() {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := rl.rotate(); err != nil {
			log.Printf("Rotation failed: %v", err)
		}
	}
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

	customLog := log.New(io.MultiWriter(os.Stdout, logger), "", log.LstdFlags)

	for i := 0; i < 100; i++ {
		customLog.Printf("Log entry %d: %s", i, time.Now().Format(time.RFC3339))
		time.Sleep(100 * time.Millisecond)
	}
}package main

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
)

const (
    maxLogSize   = 1024 * 1024 // 1MB
    maxBackups   = 5
    logFileName  = "app.log"
)

func rotateLogIfNeeded() error {
    info, err := os.Stat(logFileName)
    if os.IsNotExist(err) {
        return nil
    }
    if err != nil {
        return fmt.Errorf("failed to stat log file: %w", err)
    }

    if info.Size() < maxLogSize {
        return nil
    }

    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
    
    if err := os.Rename(logFileName, backupName); err != nil {
        return fmt.Errorf("failed to rename log file: %w", err)
    }

    backups, err := filepath.Glob(logFileName + ".*")
    if err != nil {
        return fmt.Errorf("failed to list backups: %w", err)
    }

    if len(backups) > maxBackups {
        for i := 0; i < len(backups)-maxBackups; i++ {
            if err := os.Remove(backups[i]); err != nil {
                return fmt.Errorf("failed to remove old backup %s: %w", backups[i], err)
            }
        }
    }

    return nil
}

func writeLog(message string) error {
    if err := rotateLogIfNeeded(); err != nil {
        return err
    }

    file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to open log file: %w", err)
    }
    defer file.Close()

    timestamp := time.Now().Format("2006-01-02 15:04:05")
    logEntry := fmt.Sprintf("[%s] %s\n", timestamp, message)
    
    if _, err := file.WriteString(logEntry); err != nil {
        return fmt.Errorf("failed to write log: %w", err)
    }

    return nil
}

func main() {
    for i := 1; i <= 100; i++ {
        message := fmt.Sprintf("Log entry number %d", i)
        if err := writeLog(message); err != nil {
            fmt.Printf("Error writing log: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }
    fmt.Println("Log rotation test completed")
}