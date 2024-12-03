package main

import (
    "sync"
    
)

type JobRequest struct {
    Count  int      `json:"count"`
    Visits []Visit  `json:"visits"`
}

type Visit struct {
    StoreID   string   `json:"store_id"`
    ImageURLs []string `json:"image_url"`
    VisitTime string   `json:"visit_time"`
}

type Job struct {
    JobID        int            `json:"job_id"`
    Status       string         `json:"status"`
    Errors       []ErrorDetail  `json:"error,omitempty"`
    ImageResults []ImageResult  `json:"image_results,omitempty"`
    mu           sync.Mutex
}

type ErrorDetail struct {
    StoreID string `json:"store_id"`
    Error   string `json:"error"`
}

type JobResponse struct {
    JobID int `json:"job_id"`
}

type StatusResponse struct {
    Status  string        `json:"status"`
    JobID   int           `json:"job_id"`
    Error   []ErrorDetail `json:"error,omitempty"`
    Results []ImageResult `json:"results,omitempty"`
}

var (
    jobs    = make(map[int]*Job)
    jobsMux sync.RWMutex
)

type ImageResult struct {
    URL       string `json:"url"`
    Width     int    `json:"width"`
    Height    int    `json:"height"`
    Perimeter int    `json:"perimeter"`
}