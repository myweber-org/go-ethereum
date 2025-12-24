package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
)

type ActivityLog struct {
    Timestamp time.Time
    UserID    string
    Action    string
    IPAddress string
}

var activityLogs []ActivityLog

func logActivity(userID, action, ipAddress string) {
    logEntry := ActivityLog{
        Timestamp: time.Now(),
        UserID:    userID,
        Action:    action,
        IPAddress: ipAddress,
    }
    activityLogs = append(activityLogs, logEntry)
    
    file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Printf("Failed to open log file: %v", err)
        return
    }
    defer file.Close()
    
    logLine := fmt.Sprintf("%s | %s | %s | %s\n", 
        logEntry.Timestamp.Format(time.RFC3339),
        logEntry.UserID,
        logEntry.Action,
        logEntry.IPAddress)
    
    if _, err := file.WriteString(logLine); err != nil {
        log.Printf("Failed to write log: %v", err)
    }
}

func activityHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("user")
    action := r.URL.Query().Get("action")
    
    if userID == "" || action == "" {
        http.Error(w, "Missing parameters", http.StatusBadRequest)
        return
    }
    
    ipAddress := r.RemoteAddr
    logActivity(userID, action, ipAddress)
    
    fmt.Fprintf(w, "Activity logged: %s performed '%s' from %s", userID, action, ipAddress)
}

func viewLogsHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Activity Logs:")
    fmt.Fprintln(w, "Timestamp | UserID | Action | IP Address")
    fmt.Fprintln(w, "----------------------------------------")
    
    for _, logEntry := range activityLogs {
        fmt.Fprintf(w, "%s | %s | %s | %s\n",
            logEntry.Timestamp.Format("2006-01-02 15:04:05"),
            logEntry.UserID,
            logEntry.Action,
            logEntry.IPAddress)
    }
}

func main() {
    http.HandleFunc("/log", activityHandler)
    http.HandleFunc("/logs", viewLogsHandler)
    
    port := ":8080"
    fmt.Printf("Server starting on port %s\n", port)
    fmt.Println("Endpoints:")
    fmt.Println("  /log?user=USER_ID&action=ACTION - Log user activity")
    fmt.Println("  /logs - View all activity logs")
    
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatal("Server failed: ", err)
    }
}