
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
	maxFileSize    = 10 * 1024 * 1024 // 10MB
	maxBackupFiles = 5
	logFileName    = "app.log"
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	currentPos int64
	basePath   string
}

func NewRotatingLogger(basePath string) (*RotatingLogger, error) {
	rl := &RotatingLogger{basePath: basePath}
	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	fullPath := filepath.Join(rl.basePath, logFileName)
	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}
	rl.file = file
	rl.currentPos = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentPos+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}
	n, err = rl.file.Write(p)
	if err == nil {
		rl.currentPos += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.file.Close(); err != nil {
		return err
	}

	oldPath := filepath.Join(rl.basePath, logFileName)
	for i := maxBackupFiles - 1; i >= 0; i-- {
		var source string
		if i == 0 {
			source = oldPath
		} else {
			source = filepath.Join(rl.basePath, fmt.Sprintf("%s.%d", logFileName, i))
		}
		dest := filepath.Join(rl.basePath, fmt.Sprintf("%s.%d", logFileName, i+1))

		if _, err := os.Stat(source); err == nil {
			if err := os.Rename(source, dest); err != nil {
				return err
			}
		}
	}

	if err := rl.cleanupOldFiles(); err != nil {
		return err
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) cleanupOldFiles() error {
	for i := maxBackupFiles + 1; i < 20; i++ {
		path := filepath.Join(rl.basePath, fmt.Sprintf("%s.%d", logFileName, i))
		if _, err := os.Stat(path); err == nil {
			if err := os.Remove(path); err != nil {
				return err
			}
		}
	}
	return nil
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func main() {
	logDir := "./logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatal(err)
	}

	rotator, err := NewRotatingLogger(logDir)
	if err != nil {
		log.Fatal(err)
	}
	defer rotator.Close()

	log.SetOutput(io.MultiWriter(os.Stdout, rotator))

	for i := 0; i < 100; i++ {
		log.Printf("Log entry %d at %s", i, time.Now().Format(time.RFC3339))
		time.Sleep(100 * time.Millisecond)
	}
}
package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "sync"
    "time"
)

type RotatingWriter struct {
    mu          sync.Mutex
    current     *os.File
    basePath    string
    maxSize     int64
    currentSize int64
    maxFiles    int
}

func NewRotatingWriter(basePath string, maxSize int64, maxFiles int) (*RotatingWriter, error) {
    w := &RotatingWriter{
        basePath: basePath,
        maxSize:  maxSize,
        maxFiles: maxFiles,
    }
    if err := w.openCurrent(); err != nil {
        return nil, err
    }
    return w, nil
}

func (w *RotatingWriter) openCurrent() error {
    f, err := os.OpenFile(w.basePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    stat, err := f.Stat()
    if err != nil {
        f.Close()
        return err
    }
    w.current = f
    w.currentSize = stat.Size()
    return nil
}

func (w *RotatingWriter) rotate() error {
    w.current.Close()
    timestamp := time.Now().Unix()
    for i := w.maxFiles - 1; i > 0; i-- {
        oldPath := w.basePath + "." + strconv.Itoa(i-1)
        newPath := w.basePath + "." + strconv.Itoa(i)
        if _, err := os.Stat(oldPath); err == nil {
            os.Rename(oldPath, newPath)
        }
    }
    backupPath := w.basePath + ".0"
    os.Rename(w.basePath, backupPath)
    return w.openCurrent()
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
    w.mu.Lock()
    defer w.mu.Unlock()
    if w.currentSize+int64(len(p)) > w.maxSize {
        if err := w.rotate(); err != nil {
            return 0, err
        }
    }
    n, err := w.current.Write(p)
    if err == nil {
        w.currentSize += int64(n)
    }
    return n, err
}

func (w *RotatingWriter) Close() error {
    w.mu.Lock()
    defer w.mu.Unlock()
    return w.current.Close()
}

func main() {
    writer, err := NewRotatingWriter("app.log", 1024*1024, 5)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create writer: %v\n", err)
        os.Exit(1)
    }
    defer writer.Close()
    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        writer.Write([]byte(msg))
        time.Sleep(100 * time.Millisecond)
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

const (
	maxFileSize = 10 * 1024 * 1024
	logDir      = "./logs"
)

type RotatingLogger struct {
	mu         sync.Mutex
	current    *os.File
	baseName   string
	fileSize   int64
	fileNumber int
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: baseName,
	}

	if err := rl.openNewFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openNewFile() error {
	rl.fileNumber++
	filename := filepath.Join(logDir, fmt.Sprintf("%s_%d.log", rl.baseName, rl.fileNumber))

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	if rl.current != nil {
		rl.current.Close()
		go rl.compressPreviousFile()
	}

	rl.current = file
	rl.fileSize = 0
	return nil
}

func (rl *RotatingLogger) compressPreviousFile() {
	if rl.fileNumber <= 1 {
		return
	}

	prevNum := rl.fileNumber - 1
	source := filepath.Join(logDir, fmt.Sprintf("%s_%d.log", rl.baseName, prevNum))
	dest := filepath.Join(logDir, fmt.Sprintf("%s_%d.log.gz", rl.baseName, prevNum))

	srcFile, err := os.Open(source)
	if err != nil {
		return
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return
	}
	defer destFile.Close()

	gzWriter := gzip.NewWriter(destFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, srcFile); err != nil {
		return
	}

	os.Remove(source)
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.fileSize+int64(len(p)) > maxFileSize {
		if err := rl.openNewFile(); err != nil {
			return 0, err
		}
	}

	n, err := rl.current.Write(p)
	if err == nil {
		rl.fileSize += int64(n)
	}
	return n, err
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
	logger, err := NewRotatingLogger("app")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(message))
		time.Sleep(50 * time.Millisecond)
	}

	fmt.Println("Log rotation completed")
}