package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"kiranaC_BE/models"
	"kiranaC_BE/utils"
)

var jobMux sync.Mutex
var jobs = make(map[int]*models.Job)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	jobIDStr := r.URL.Query().Get("jobid")
	if jobIDStr == "" {
		http.Error(w, "JobID is required", http.StatusBadRequest)
		return
	}

	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		http.Error(w, "Invalid JobID format", http.StatusBadRequest)
		return
	}

	utils.JobMux.Lock()
	job, exists := utils.Jobs[jobID]
	utils.JobMux.Unlock()

	if !exists {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}
