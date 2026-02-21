
package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize    = 10 * 1024 * 1024 // 10MB
    maxBackupFiles = 5
    logFileName    = "app.log"
)

type LogRotator struct {
    currentFile *os.File
    filePath    string
    bytesWritten int64
}

func NewLogRotator() (*LogRotator, error) {
    lr := &LogRotator{
        filePath: logFileName,
    }
    
    if err := lr.openCurrentFile(); err != nil {
        return nil, err
    }
    
    return lr, nil
}

func (lr *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    
    lr.currentFile = file
    lr.bytesWritten = info.Size()
    
    return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.bytesWritten+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }
    
    n, err := lr.currentFile.Write(p)
    if err != nil {
        return n, err
    }
    
    lr.bytesWritten += int64(n)
    return n, nil
}

func (lr *LogRotator) rotate() error {
    if err := lr.currentFile.Close(); err != nil {
        return err
    }
    
    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s", lr.filePath, timestamp)
    
    if err := os.Rename(lr.filePath, backupName); err != nil {
        return err
    }
    
    if err := lr.openCurrentFile(); err != nil {
        return err
    }
    
    lr.cleanupOldBackups()
    
    return nil
}

func (lr *LogRotator) cleanupOldBackups() {
    pattern := fmt.Sprintf("%s.*", lr.filePath)
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }
    
    if len(matches) <= maxBackupFiles {
        return
    }
    
    sort.Sort(sort.Reverse(sort.StringSlice(matches)))
    
    for i := maxBackupFiles; i < len(matches); i++ {
        os.Remove(matches[i])
    }
}

func (lr *LogRotator) Close() error {
    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func extractTimestamp(filename string) (time.Time, error) {
    parts := strings.Split(filename, ".")
    if len(parts) < 2 {
        return time.Time{}, fmt.Errorf("invalid backup filename")
    }
    
    timestampStr := parts[len(parts)-1]
    return time.Parse("20060102_150405", timestampStr)
}

func main() {
    rotator, err := NewLogRotator()
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        os.Exit(1)
    }
    defer rotator.Close()
    
    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
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
    sequence    int
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        basePath: path,
    }
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
    if rl.currentFile != nil {
        rl.currentFile.Close()
        if err := rl.compressCurrentFile(); err != nil {
            return err
        }
    }
    rl.sequence++
    if rl.sequence > maxBackups {
        rl.cleanOldBackups()
    }
    return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressCurrentFile() error {
    oldPath := rl.getFilePath(rl.sequence)
    compressedPath := oldPath + ".gz"
    
    src, err := os.Open(oldPath)
    if err != nil {
        return err
    }
    defer src.Close()
    
    dst, err := os.Create(compressedPath)
    if err != nil {
        return err
    }
    defer dst.Close()
    
    gz := gzip.NewWriter(dst)
    defer gz.Close()
    
    _, err = io.Copy(gz, src)
    if err != nil {
        return err
    }
    
    return os.Remove(oldPath)
}

func (rl *RotatingLogger) openCurrentFile() error {
    path := rl.getFilePath(0)
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }
    
    f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

func (rl *RotatingLogger) getFilePath(seq int) string {
    if seq == 0 {
        return rl.basePath
    }
    timestamp := time.Now().Format("20060102_150405")
    return fmt.Sprintf("%s.%d.%s", rl.basePath, seq, timestamp)
}

func (rl *RotatingLogger) cleanOldBackups() {
    for i := rl.sequence - maxBackups; i > 0; i-- {
        pattern := fmt.Sprintf("%s.%d.*.gz", rl.basePath, i)
        matches, _ := filepath.Glob(pattern)
        for _, match := range matches {
            os.Remove(match)
        }
    }
}

func (rl *RotatingLogger) Close() error {
    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("./logs/app.log")
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
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	maxLogSize   = 1024 * 1024 // 1MB
	maxBackupLog = 5
	logFileName  = "app.log"
)

type LogRotator struct {
	currentSize int64
	file        *os.File
	logger      *log.Logger
}

func NewLogRotator() (*LogRotator, error) {
	lr := &LogRotator{}
	if err := lr.openLogFile(); err != nil {
		return nil, err
	}
	lr.logger = log.New(lr.file, "", log.LstdFlags)
	return lr, nil
}

func (lr *LogRotator) openLogFile() error {
	info, err := os.Stat(logFileName)
	if os.IsNotExist(err) {
		file, err := os.Create(logFileName)
		if err != nil {
			return err
		}
		lr.file = file
		lr.currentSize = 0
		return nil
	}
	if err != nil {
		return err
	}

	lr.currentSize = info.Size()
	if lr.currentSize >= maxLogSize {
		if err := lr.rotate(); err != nil {
			return err
		}
		file, err := os.Create(logFileName)
		if err != nil {
			return err
		}
		lr.file = file
		lr.currentSize = 0
	} else {
		file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		lr.file = file
	}
	return nil
}

func (lr *LogRotator) rotate() error {
	if err := lr.file.Close(); err != nil {
		return err
	}

	for i := maxBackupLog - 1; i >= 0; i-- {
		oldName := fmt.Sprintf("%s.%d", logFileName, i)
		newName := fmt.Sprintf("%s.%d", logFileName, i+1)
		if _, err := os.Stat(oldName); err == nil {
			if err := os.Rename(oldName, newName); err != nil {
				return err
			}
		}
	}

	backupName := fmt.Sprintf("%s.0", logFileName)
	if err := os.Rename(logFileName, backupName); err != nil {
		return err
	}

	oldestBackup := fmt.Sprintf("%s.%d", logFileName, maxBackupLog)
	if _, err := os.Stat(oldestBackup); err == nil {
		if err := os.Remove(oldestBackup); err != nil {
			return err
		}
	}
	return nil
}

func (lr *LogRotator) Write(p []byte) (n int, err error) {
	if lr.currentSize+int64(len(p)) >= maxLogSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
		if err := lr.openLogFile(); err != nil {
			return 0, err
		}
	}
	n, err = lr.file.Write(p)
	if err == nil {
		lr.currentSize += int64(n)
	}
	return n, err
}

func (lr *LogRotator) Println(v ...interface{}) {
	message := fmt.Sprintln(v...)
	lr.Write([]byte(message))
}

func (lr *LogRotator) Cleanup() error {
	if lr.file != nil {
		return lr.file.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator()
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		return
	}
	defer rotator.Cleanup()

	for i := 0; i < 100; i++ {
		rotator.Println("Log entry number:", i, "at", time.Now().Format(time.RFC3339))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed. Check app.log and backup files.")
}