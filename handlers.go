package main

import (
    "encoding/json"
    "net/http"
    "strconv"
    "time"
	"log"
)

func SubmitJobHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received submit request")
    var req JobRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        log.Printf("Error decoding request: %v", err)
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
        return
    }

    if req.Count != len(req.Visits) {
		log.Printf("Count mismatch: expected %d, got %d visits", req.Count, len(req.Visits))
        http.Error(w, `{"error": "Count does not match number of visits"}`, http.StatusBadRequest)
        return
    }
	// Validate store IDs first
    for _, visit := range req.Visits {
        if !ValidateStore(visit.StoreID) {
            log.Printf("Invalid store ID: %s", visit.StoreID)
            http.Error(w, `{"error": "Invalid store ID"}`, http.StatusBadRequest)
            return
        }
    }

    // Create new job
    jobID := int(time.Now().UnixNano() % 1000000)
    job := &Job{
        JobID:  jobID,
        Status: "ongoing",
    }

    // Store job
    jobsMux.Lock()
    jobs[jobID] = job
    jobsMux.Unlock()

    // Process job in background
    go ProcessJob(job, req.Visits)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    response := JobResponse{JobID: jobID}
    json.NewEncoder(w).Encode(response)
    log.Printf("Created job with ID: %d", jobID)
}

func GetJobStatusHandler(w http.ResponseWriter, r *http.Request) {
    jobID, err := strconv.Atoi(r.URL.Query().Get("jobid"))
    if err != nil {
        http.Error(w, "{}", http.StatusBadRequest)
        return
    }

    jobsMux.RLock()
    job, exists := jobs[jobID]
    jobsMux.RUnlock()

    if !exists {
        http.Error(w, "{}", http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(StatusResponse{
        Status: job.Status,
        JobID:  job.JobID,
        Error:  job.Errors,
    })
}