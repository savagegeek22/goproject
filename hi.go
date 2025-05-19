// File: main.go
package main

import (
    "encoding/json"
    "log"
    "net/http"
)

type Message struct {
    Greeting string `json:"greeting"`
}

func greetHandler(w http.ResponseWriter, r *http.Request) {
    msg := Message{Greeting: "Hello, Golang Developer!"}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(msg)
}

func main() {
    http.HandleFunc("/greet", greetHandler)
    log.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

