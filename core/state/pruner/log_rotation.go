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
    mu         sync.Mutex
    file       *os.File
    currentSize int64
    basePath   string
    sequence   int
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
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

    return &RotatingLogger{
        file:       file,
        currentSize: info.Size(),
        basePath:   path,
        sequence:   0,
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > maxFileSize {
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
    rotatedPath := fmt.Sprintf("%s.%s.%d", rl.basePath, timestamp, rl.sequence)
    if err := os.Rename(rl.basePath, rotatedPath); err != nil {
        return err
    }

    if err := rl.compressFile(rotatedPath); err != nil {
        return err
    }

    rl.sequence++
    if rl.sequence > maxBackups {
        rl.cleanOldBackups()
    }

    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.file = file
    rl.currentSize = 0
    return nil
}

func (rl *RotatingLogger) compressFile(src string) error {
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

    gz := gzip.NewWriter(dstFile)
    defer gz.Close()

    if _, err := io.Copy(gz, srcFile); err != nil {
        return err
    }

    return os.Remove(src)
}

func (rl *RotatingLogger) cleanOldBackups() {
    pattern := rl.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) > maxBackups {
        toDelete := matches[:len(matches)-maxBackups]
        for _, path := range toDelete {
            os.Remove(path)
        }
    }
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.file.Close()
}

func main() {
    logger, err := NewRotatingLogger("logs/application.log")
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }
}
package main

import (
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
    logger := &RotatingLogger{
        filename: filename,
    }
    
    if err := logger.openCurrent(); err != nil {
        return nil, err
    }
    
    return logger, nil
}

func (l *RotatingLogger) openCurrent() error {
    file, err := os.OpenFile(l.filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    
    l.current = file
    l.size = info.Size()
    return nil
}

func (l *RotatingLogger) Write(p []byte) (int, error) {
    l.mu.Lock()
    defer l.mu.Unlock()
    
    if l.size+int64(len(p)) > maxFileSize {
        if err := l.rotate(); err != nil {
            return 0, err
        }
    }
    
    n, err := l.current.Write(p)
    if err == nil {
        l.size += int64(n)
    }
    return n, err
}

func (l *RotatingLogger) rotate() error {
    if l.current != nil {
        l.current.Close()
    }
    
    for i := maxBackups - 1; i >= 0; i-- {
        oldName := l.backupName(i)
        newName := l.backupName(i + 1)
        
        if _, err := os.Stat(oldName); err == nil {
            os.Rename(oldName, newName)
        }
    }
    
    if err := os.Rename(l.filename, l.backupName(0)); err != nil && !os.IsNotExist(err) {
        return err
    }
    
    return l.openCurrent()
}

func (l *RotatingLogger) backupName(index int) string {
    if index == 0 {
        return l.filename + ".1"
    }
    return fmt.Sprintf("%s.%d", l.filename, index+1)
}

func (l *RotatingLogger) Close() error {
    l.mu.Lock()
    defer l.mu.Unlock()
    
    if l.current != nil {
        return l.current.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app.log")
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()
    
    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("[%s] Log entry %d: Application is running normally\n", 
            time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(message))
        time.Sleep(10 * time.Millisecond)
    }
    
    fmt.Println("Log rotation test completed")
}