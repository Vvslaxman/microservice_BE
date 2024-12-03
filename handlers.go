package main

import (
    "encoding/json"
    "net/http"
    "strconv"
    "time"
	"log"
	"github.com/gorilla/mux"
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

    // Validate store IDs and image URLs
    for _, visit := range req.Visits {
        // First, validate store ID
        if !ValidateStore(visit.StoreID) {
            log.Printf("Invalid store ID: %s", visit.StoreID)
            http.Error(w, `{"error": "Invalid store ID"}`, http.StatusBadRequest)
            return
        }

        // Check if image URLs exist and are not empty
        if len(visit.ImageURLs) == 0 {
            log.Printf("No image URLs for store ID: %s", visit.StoreID)
            http.Error(w, `{"error": "No image URLs provided"}`, http.StatusBadRequest)
            return
        }

        // Additional check for empty URLs
        for _, url := range visit.ImageURLs {
            if url == "" {
                log.Printf("Empty image URL for store ID: %s", visit.StoreID)
                http.Error(w, `{"error": "Empty image URL found"}`, http.StatusBadRequest)
                return
            }
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
    vars := mux.Vars(r)
    jobIDStr := r.URL.Query().Get("jobid")
    
    // Fallback to path variable if query parameter is not present
    if jobIDStr == "" {
        jobIDStr = vars["jobid"]
    }

    jobID, err := strconv.Atoi(jobIDStr)
    if err != nil {
        log.Printf("Invalid job ID: %s", jobIDStr)
        http.Error(w, `{"error": "Invalid job ID"}`, http.StatusBadRequest)
        return
    }

    jobsMux.RLock()
    job, exists := jobs[jobID]
    jobsMux.RUnlock()

    if !exists {
        log.Printf("Job not found: %d", jobID)
        http.Error(w, `{"error": "Job not found"}`, http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(StatusResponse{
        Status:  job.Status,
        JobID:   job.JobID,
        Error:   job.Errors,
        Results: job.ImageResults,
    })
}