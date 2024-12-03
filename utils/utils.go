package utils

import (
		
	"kiranaC_BE/models"
	"sync"
	"time"
	
)

// Exported variables
var JobMux sync.Mutex
var Jobs = make(map[int]*models.Job)

// CreateJob creates a new job based on the job request data.
func CreateJob(request models.JobRequest) int {
    // Generate a new job ID (could be more sophisticated)
    jobID := int(time.Now().UnixNano())

    // Convert []string to []Visit
    visits := make([]models.Visit, len(request.Visits))
    for i, visit := range request.Visits {
        visits[i] = models.Visit{StoreID: visit}
    }

    // Create a new job object
    job := &models.Job{
        JobID:  jobID,
        Status: "ongoing",
        Visits: visits,
    }

    // Store the job in a global jobs map (you might need synchronization for thread safety)
    JobMux.Lock()
    Jobs[jobID] = job
    JobMux.Unlock()

    return jobID
}
