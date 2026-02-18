package main

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "sync"
)

type FileStats struct {
    Path     string
    Size     int64
    Lines    int
    Error    error
}

func processFile(path string, results chan<- FileStats, wg *sync.WaitGroup) {
    defer wg.Done()
    
    stats := FileStats{Path: path}
    
    file, err := os.Open(path)
    if err != nil {
        stats.Error = err
        results <- stats
        return
    }
    defer file.Close()
    
    info, err := file.Stat()
    if err != nil {
        stats.Error = err
        results <- stats
        return
    }
    stats.Size = info.Size()
    
    scanner := bufio.NewScanner(file)
    lineCount := 0
    for scanner.Scan() {
        lineCount++
    }
    
    if err := scanner.Err(); err != nil {
        stats.Error = err
    } else {
        stats.Lines = lineCount
    }
    
    results <- stats
}

func collectStats(paths []string) []FileStats {
    var wg sync.WaitGroup
    results := make(chan FileStats, len(paths))
    
    for _, path := range paths {
        wg.Add(1)
        go processFile(path, results, &wg)
    }
    
    wg.Wait()
    close(results)
    
    var allStats []FileStats
    for stat := range results {
        allStats = append(allStats, stat)
    }
    
    return allStats
}

func findFiles(dir string, pattern string) ([]string, error) {
    var files []string
    
    err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if !info.IsDir() {
            matched, err := filepath.Match(pattern, filepath.Base(path))
            if err != nil {
                return err
            }
            if matched {
                files = append(files, path)
            }
        }
        
        return nil
    })
    
    return files, err
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: file_processor <directory> <pattern>")
        fmt.Println("Example: file_processor . *.txt")
        os.Exit(1)
    }
    
    dir := os.Args[1]
    pattern := os.Args[2]
    
    files, err := findFiles(dir, pattern)
    if err != nil {
        fmt.Printf("Error finding files: %v\n", err)
        os.Exit(1)
    }
    
    if len(files) == 0 {
        fmt.Println("No files found matching pattern")
        return
    }
    
    fmt.Printf("Processing %d files...\n", len(files))
    
    stats := collectStats(files)
    
    var totalSize int64
    var totalLines int
    var errorCount int
    
    for _, stat := range stats {
        if stat.Error != nil {
            fmt.Printf("Error processing %s: %v\n", stat.Path, stat.Error)
            errorCount++
            continue
        }
        
        totalSize += stat.Size
        totalLines += stat.Lines
        fmt.Printf("%s: %d bytes, %d lines\n", stat.Path, stat.Size, stat.Lines)
    }
    
    fmt.Printf("\nSummary:\n")
    fmt.Printf("Total files processed: %d\n", len(files))
    fmt.Printf("Successful: %d\n", len(files)-errorCount)
    fmt.Printf("Errors: %d\n", errorCount)
    fmt.Printf("Total size: %d bytes\n", totalSize)
    fmt.Printf("Total lines: %d\n", totalLines)
}