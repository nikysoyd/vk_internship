package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"

    _ "github.com/lib/pq"
)

type PingResult struct {
    IP          string    `json:"ip"`
    PingTime    time.Time `json:"ping_time"`
    IsSuccess   bool      `json:"is_success"`
    LastSuccess time.Time `json:"last_success"`
}

var db *sql.DB

func initDB() {
    var err error
    connStr := "user=postgres dbname=docker_pinger sslmode=disable password=postgres host=postgres"
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
}

func getPingResults(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT ip, ping_time, is_success, last_success FROM ping_results")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var results []PingResult
    for rows.Next() {
        var res PingResult
        if err := rows.Scan(&res.IP, &res.PingTime, &res.IsSuccess, &res.LastSuccess); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        results = append(results, res)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(results)
}

func addPingResult(w http.ResponseWriter, r *http.Request) {
    var res PingResult
    if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    _, err := db.Exec("INSERT INTO ping_results (ip, ping_time, is_success, last_success) VALUES ($1, $2, $3, $4)",
        res.IP, res.PingTime, res.IsSuccess, res.LastSuccess)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func main() {
    initDB()
    defer db.Close()

    http.HandleFunc("/ping-results", getPingResults)
    http.HandleFunc("/add-ping-result", addPingResult)

    fmt.Println("Server is running on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
