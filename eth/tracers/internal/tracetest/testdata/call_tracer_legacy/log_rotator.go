
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
	maxFileSize = 10 * 1024 * 1024
	backupCount = 5
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	size       int64
	basePath   string
	currentNum int
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: path,
	}
	if err := rl.openCurrent(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrent() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		rl.file.Close()
	}

	file, err := os.OpenFile(rl.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.file = file
	rl.size = stat.Size()
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

	n, err := rl.file.Write(p)
	if err == nil {
		rl.size += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.file.Close(); err != nil {
		return err
	}

	for i := backupCount - 1; i >= 0; i-- {
		oldPath := rl.getBackupPath(i)
		newPath := rl.getBackupPath(i + 1)

		if _, err := os.Stat(oldPath); err == nil {
			if i == backupCount-1 {
				os.Remove(oldPath)
			} else {
				if err := os.Rename(oldPath, newPath); err != nil {
					return err
				}
			}
		}
	}

	firstBackup := rl.getBackupPath(0)
	if err := os.Rename(rl.basePath, firstBackup); err != nil {
		return err
	}

	if err := rl.compressFile(firstBackup); err != nil {
		return err
	}

	return rl.openCurrent()
}

func (rl *RotatingLogger) getBackupPath(num int) string {
	if num == 0 {
		return rl.basePath + ".1"
	}
	return fmt.Sprintf("%s.%d.gz", rl.basePath, num+1)
}

func (rl *RotatingLogger) compressFile(src string) error {
	dest := src + ".gz"

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	gzWriter := gzip.NewWriter(destFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, srcFile); err != nil {
		return err
	}

	os.Remove(src)
	return nil
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		return rl.file.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n",
			time.Now().Format("2006-01-02 15:04:05"), i)
		logger.Write([]byte(msg))
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

type RotatingLog struct {
    basePath      string
    maxSize       int64
    currentSize   int64
    currentFile   *os.File
    fileCounter   int
    compressOld   bool
}

func NewRotatingLog(basePath string, maxSize int64, compressOld bool) (*RotatingLog, error) {
    rl := &RotatingLog{
        basePath:    basePath,
        maxSize:     maxSize,
        compressOld: compressOld,
    }

    err := rl.initialize()
    if err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLog) initialize() error {
    dir := filepath.Dir(rl.basePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }

    existingFiles, err := filepath.Glob(rl.basePath + ".*")
    if err != nil {
        return err
    }

    maxNum := 0
    for _, f := range existingFiles {
        suffix := strings.TrimPrefix(f, rl.basePath+".")
        if num, err := strconv.Atoi(suffix); err == nil && num > maxNum {
            maxNum = num
        }
    }
    rl.fileCounter = maxNum

    currentPath := rl.basePath
    if maxNum > 0 {
        currentPath = fmt.Sprintf("%s.%d", rl.basePath, maxNum)
    }

    file, err := os.OpenFile(currentPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func (rl *RotatingLog) Write(p []byte) (int, error) {
    if rl.currentSize+int64(len(p)) > rl.maxSize && rl.currentSize > 0 {
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

func (rl *RotatingLog) rotate() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }

    rl.fileCounter++
    newPath := fmt.Sprintf("%s.%d", rl.basePath, rl.fileCounter)

    if err := os.Rename(rl.basePath, newPath); err != nil {
        return err
    }

    if rl.compressOld {
        go rl.compressFile(newPath)
    }

    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.currentFile = file
    rl.currentSize = 0

    return nil
}

func (rl *RotatingLog) compressFile(path string) error {
    compressedPath := path + ".gz"

    src, err := os.Open(path)
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

    if _, err := io.Copy(gz, src); err != nil {
        return err
    }

    os.Remove(path)
    return nil
}

func (rl *RotatingLog) Close() error {
    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    log, err := NewRotatingLog("/var/log/myapp/app.log", 1024*1024, true)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer log.Close()

    for i := 0; i < 100; i++ {
        message := fmt.Sprintf("[%s] Log entry %d: Application is running normally\n",
            time.Now().Format(time.RFC3339), i)
        log.Write([]byte(message))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}