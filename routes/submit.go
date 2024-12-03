package routes

import (
	"encoding/json"
	"kiranaC_BE/models"
	"kiranaC_BE/utils"
	"net/http"
	"log"
)

// SubmitJobHandler handles job submission requests
func SubmitJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var jobRequest models.JobRequest
	if err := json.NewDecoder(r.Body).Decode(&jobRequest); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Validate the job request (you can expand this as needed)
	if len(jobRequest.Visits) == 0 {
		http.Error(w, "No visits provided", http.StatusBadRequest)
		return
	}

	// Create a job and return the job ID
	jobID := utils.CreateJob(jobRequest)
	response := map[string]interface{}{
		"job_id": jobID,
		"status": "Job created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Println("Received request at /api/submit")
	log.Printf("Method: %s, URL: %s", r.Method, r.URL.Path)
}
