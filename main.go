package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func main() {
    // Load store data
    if err := LoadStoreData("data/StoreMasterAssignment.csv"); err != nil {
        log.Fatalf("Failed to load store data: %v", err)
    }

    r := mux.NewRouter()
    r.HandleFunc("/api/submit/", SubmitJobHandler).Methods("POST")
    r.HandleFunc("/api/status", GetJobStatusHandler).Methods("GET")

    log.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
