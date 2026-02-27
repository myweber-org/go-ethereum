
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
	filePath     string
	maxSize      int64
	currentSize  int64
	rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	rl := &RotatingLogger{
		filePath: basePath,
		maxSize:  maxSize,
	}
	
	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}
	
	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	timestamp := time.Now().Format("20060102_150405")
	archivePath := fmt.Sprintf("%s.%s.%d.gz", rl.filePath, timestamp, rl.rotationCount)
	
	if err := rl.compressFile(rl.filePath, archivePath); err != nil {
		return err
	}
	
	if err := os.Remove(rl.filePath); err != nil && !os.IsNotExist(err) {
		return err
	}
	
	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(source, destination string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	
	destFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer destFile.Close()
	
	gzWriter := gzip.NewWriter(destFile)
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
	logger, err := NewRotatingLogger("app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()
	
	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: This is a sample log message\n", 
			time.Now().Format(time.RFC3339), i)
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
	
	fmt.Println("Log rotation test completed")
}package main

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
    maxFileSize = 10 * 1024 * 1024
    maxBackups  = 5
)

type RotatingLog struct {
    filePath string
    file     *os.File
    size     int64
}

func NewRotatingLog(path string) (*RotatingLog, error) {
    rl := &RotatingLog{filePath: path}
    if err := rl.openCurrent(); err != nil {
        return nil, err
    }
    return rl, nil
}

func (rl *RotatingLog) Write(p []byte) (int, error) {
    if rl.size+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.file.Write(p)
    if err == nil {
        rl.size += int64(n)
    }
    return n, err
}

func (rl *RotatingLog) openCurrent() error {
    info, err := os.Stat(rl.filePath)
    if err != nil && !os.IsNotExist(err) {
        return err
    }

    if err == nil {
        rl.size = info.Size()
    } else {
        rl.size = 0
    }

    rl.file, err = os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    return err
}

func (rl *RotatingLog) rotate() error {
    if rl.file != nil {
        rl.file.Close()
    }

    timestamp := time.Now().Format("20060102150405")
    rotatedPath := fmt.Sprintf("%s.%s", rl.filePath, timestamp)

    if err := os.Rename(rl.filePath, rotatedPath); err != nil {
        return err
    }

    if err := rl.compressFile(rotatedPath); err != nil {
        return err
    }

    if err := rl.cleanupOldBackups(); err != nil {
        return err
    }

    return rl.openCurrent()
}

func (rl *RotatingLog) compressFile(src string) error {
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    dstFile, err := os.Create(src + ".gz")
    if err != nil {
        return err
    }
    defer dstFile.Close()

    gzWriter := gzip.NewWriter(dstFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, srcFile); err != nil {
        return err
    }

    return os.Remove(src)
}

func (rl *RotatingLog) cleanupOldBackups() error {
    pattern := rl.filePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= maxBackups {
        return nil
    }

    backupMap := make(map[int64]string)
    for _, match := range matches {
        parts := strings.Split(filepath.Base(match), ".")
        if len(parts) < 3 {
            continue
        }
        timestamp, err := strconv.ParseInt(parts[len(parts)-2], 10, 64)
        if err != nil {
            continue
        }
        backupMap[timestamp] = match
    }

    timestamps := make([]int64, 0, len(backupMap))
    for ts := range backupMap {
        timestamps = append(timestamps, ts)
    }

    for i := 0; i < len(timestamps)-maxBackups; i++ {
        oldest := timestamps[i]
        if err := os.Remove(backupMap[oldest]); err != nil {
            return err
        }
    }

    return nil
}

func (rl *RotatingLog) Close() error {
    if rl.file != nil {
        return rl.file.Close()
    }
    return nil
}

func main() {
    log, err := NewRotatingLog("application.log")
    if err != nil {
        fmt.Printf("Failed to create log: %v\n", err)
        return
    }
    defer log.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := log.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}