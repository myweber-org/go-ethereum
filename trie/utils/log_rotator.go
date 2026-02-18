
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
    maxFileSize   = 10 * 1024 * 1024 // 10MB
    maxBackupFiles = 5
)

type LogRotator struct {
    currentFile *os.File
    currentSize int64
    basePath    string
}

func NewLogRotator(basePath string) (*LogRotator, error) {
    rotator := &LogRotator{basePath: basePath}
    if err := rotator.openCurrentFile(); err != nil {
        return nil, err
    }
    return rotator, nil
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
    if lr.currentFile != nil {
        lr.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", lr.basePath, timestamp)
    if err := os.Rename(lr.basePath, rotatedPath); err != nil {
        return err
    }

    if err := lr.compressFile(rotatedPath); err != nil {
        return err
    }

    lr.cleanupOldFiles()
    return lr.openCurrentFile()
}

func (lr *LogRotator) compressFile(sourcePath string) error {
    sourceFile, err := os.Open(sourcePath)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    compressedPath := sourcePath + ".gz"
    compressedFile, err := os.Create(compressedPath)
    if err != nil {
        return err
    }
    defer compressedFile.Close()

    gzWriter := gzip.NewWriter(compressedFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, sourceFile); err != nil {
        return err
    }

    os.Remove(sourcePath)
    return nil
}

func (lr *LogRotator) cleanupOldFiles() {
    pattern := lr.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackupFiles {
        return
    }

    var timestamps []time.Time
    fileMap := make(map[time.Time]string)

    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 3 {
            continue
        }
        timestampStr := parts[len(parts)-2]
        t, err := time.Parse("20060102_150405", timestampStr)
        if err != nil {
            continue
        }
        timestamps = append(timestamps, t)
        fileMap[t] = match
    }

    for i := 0; i < len(timestamps)-maxBackupFiles; i++ {
        oldest := timestamps[i]
        if path, exists := fileMap[oldest]; exists {
            os.Remove(path)
        }
    }
}

func (lr *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(lr.basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    lr.currentFile = file
    lr.currentSize = stat.Size()
    return nil
}

func (lr *LogRotator) Close() error {
    if lr.currentFile != nil {
        return lr.currentFile.Close()
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

    fmt.Println("Log rotation test completed")
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

type LogRotator struct {
	mu          sync.Mutex
	currentFile *os.File
	filePath    string
	maxSize     int64
	backupCount int
	bytesWritten int64
}

func NewLogRotator(filePath string, maxSizeMB int, backupCount int) (*LogRotator, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	
	rotator := &LogRotator{
		filePath:    filePath,
		maxSize:     maxSize,
		backupCount: backupCount,
	}
	
	if err := rotator.openCurrentFile(); err != nil {
		return nil, err
	}
	
	return rotator, nil
}

func (r *LogRotator) openCurrentFile() error {
	dir := filepath.Dir(r.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directory failed: %v", err)
	}
	
	file, err := os.OpenFile(r.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("open file failed: %v", err)
	}
	
	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return fmt.Errorf("stat file failed: %v", err)
	}
	
	r.currentFile = file
	r.bytesWritten = stat.Size()
	return nil
}

func (r *LogRotator) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if r.bytesWritten+int64(len(p)) > r.maxSize {
		if err := r.rotate(); err != nil {
			return 0, err
		}
	}
	
	n, err := r.currentFile.Write(p)
	if err != nil {
		return n, err
	}
	
	r.bytesWritten += int64(n)
	return n, nil
}

func (r *LogRotator) rotate() error {
	if r.currentFile != nil {
		r.currentFile.Close()
	}
	
	if err := r.compressOldLogs(); err != nil {
		return fmt.Errorf("compress old logs failed: %v", err)
	}
	
	if err := r.renameCurrentFile(); err != nil {
		return fmt.Errorf("rename current file failed: %v", err)
	}
	
	if err := r.cleanupOldBackups(); err != nil {
		return fmt.Errorf("cleanup old backups failed: %v", err)
	}
	
	return r.openCurrentFile()
}

func (r *LogRotator) compressOldLogs() error {
	for i := r.backupCount - 1; i >= 0; i-- {
		sourcePath := r.getBackupPath(i)
		compressedPath := sourcePath + ".gz"
		
		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			continue
		}
		
		if _, err := os.Stat(compressedPath); err == nil {
			continue
		}
		
		if err := compressFile(sourcePath, compressedPath); err != nil {
			return err
		}
		
		os.Remove(sourcePath)
	}
	return nil
}

func compressFile(source, destination string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	
	dstFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	
	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()
	
	_, err = io.Copy(gzWriter, srcFile)
	return err
}

func (r *LogRotator) renameCurrentFile() error {
	timestamp := time.Now().Format("20060102_150405")
	newPath := fmt.Sprintf("%s.%s", r.filePath, timestamp)
	return os.Rename(r.filePath, newPath)
}

func (r *LogRotator) getBackupPath(index int) string {
	timestamp := time.Now().Add(time.Duration(-index) * 24 * time.Hour).Format("20060102")
	return fmt.Sprintf("%s.%s", r.filePath, timestamp)
}

func (r *LogRotator) cleanupOldBackups() error {
	pattern := r.filePath + ".*"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	
	if len(matches) <= r.backupCount {
		return nil
	}
	
	for i := r.backupCount; i < len(matches); i++ {
		os.Remove(matches[i])
	}
	
	return nil
}

func (r *LogRotator) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if r.currentFile != nil {
		return r.currentFile.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator("/var/log/myapp/app.log", 10, 5)
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		return
	}
	defer rotator.Close()
	
	for i := 0; i < 1000; i++ {
		logEntry := fmt.Sprintf("[%s] Log entry %d: Application is running normally\n", 
			time.Now().Format(time.RFC3339), i)
		if _, err := rotator.Write([]byte(logEntry)); err != nil {
			fmt.Printf("Write failed: %v\n", err)
			break
		}
		
		time.Sleep(10 * time.Millisecond)
	}
	
	fmt.Println("Log rotation test completed")
}