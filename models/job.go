package models

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type SubmitJobRequest struct {
	Count  int     `json:"count"`
	Visits []Visit `json:"visits"`
}

type Visit struct {
	StoreID   string   `json:"store_id"`
	ImageURLs []string `json:"image_url"`
	VisitTime string   `json:"visit_time"`
}
type JobRequest struct {
    Visits []string `json:"visits"`
}

type Job struct {
	JobID   int
	Status  string
	Visits  []Visit
	Results []string
	Errors  []JobError
	mu      sync.Mutex
}

type JobError struct {
	StoreID string `json:"store_id"`
	Error   string `json:"error"`
}

func NewJob(jobID int, visits []Visit) *Job {
	return &Job{
		JobID:   jobID,
		Status:  "ongoing",
		Visits:  visits,
		Results: []string{},
		Errors:  []JobError{},
	}
}

// ProcessJob processes all visits in a job and validates store IDs
func (j *Job) ProcessJob(storeData map[string]*Store) {
	for _, visit := range j.Visits {
		if _, valid := storeData[visit.StoreID]; !valid {
			j.mu.Lock()
			j.Errors = append(j.Errors, JobError{StoreID: visit.StoreID, Error: "Invalid StoreID"})
			j.mu.Unlock()
			continue
		}
		for _, url := range visit.ImageURLs {
			time.Sleep(time.Duration(rand.Intn(300)+100) * time.Millisecond)
			_, err := http.Get(url)

			j.mu.Lock()
			if err != nil {
				j.Errors = append(j.Errors, JobError{StoreID: visit.StoreID, Error: err.Error()})
			} else {
				result := fmt.Sprintf("Image processed for store %s: [URL: %s]", visit.StoreID, url)
				j.Results = append(j.Results, result)
			}
			j.mu.Unlock()
		}
	}

	j.mu.Lock()
	if len(j.Errors) > 0 {
		j.Status = "failed"
	} else {
		j.Status = "completed"
	}
	j.mu.Unlock()
}
