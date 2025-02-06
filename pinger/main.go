package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os/exec"
    "time"
)

type PingResult struct {
    IP          string    `json:"ip"`
    PingTime    time.Time `json:"ping_time"`
    IsSuccess   bool      `json:"is_success"`
    LastSuccess time.Time `json:"last_success"`
}

func pingIP(ip string) bool {
    cmd := exec.Command("ping", "-c", "1", ip)
    err := cmd.Run()
    return err == nil
}

func sendPingResult(result PingResult) {
    jsonData, err := json.Marshal(result)
    if err != nil {
        fmt.Println("Error marshalling JSON:", err)
        return
    }

    resp, err := http.Post("http://backend:8080/add-ping-result", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Println("Error sending data to backend:", err)
        return
    }
    defer resp.Body.Close()
}

func main() {
    ips := []string{"172.17.0.2", "172.17.0.3"} // Пример IP-адресов контейнеров

    for {
        for _, ip := range ips {
            isSuccess := pingIP(ip)
            lastSuccess := time.Time{}
            if isSuccess {
                lastSuccess = time.Now()
            }

            result := PingResult{
                IP:          ip,
                PingTime:    time.Now(),
                IsSuccess:   isSuccess,
                LastSuccess: lastSuccess,
            }

            sendPingResult(result)
        }

        time.Sleep(10 * time.Second) // Интервал пинга
    }
}
