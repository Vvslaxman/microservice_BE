package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func ProcessJob(job *Job, visits []Visit) {
	var wg sync.WaitGroup

	for _, visit := range visits {
		// Validate store
		if !ValidateStore(visit.StoreID) {
			job.mu.Lock()
			job.Errors = append(job.Errors, ErrorDetail{
				StoreID: visit.StoreID,
				Error:   "Invalid store ID",
			})
			job.mu.Unlock()
			continue
		}

		for _, url := range visit.ImageURLs {
			wg.Add(1)
			go func(storeID, imageURL string) {
				defer wg.Done()

				result, err := processImage(imageURL)
				if err != nil {
					job.mu.Lock()
					job.Errors = append(job.Errors, ErrorDetail{
						StoreID: storeID,
						Error:   err.Error(),
					})
					job.mu.Unlock()
				} else {
					job.mu.Lock()
					job.ImageResults = append(job.ImageResults, result)
					job.mu.Unlock()
				}
			}(visit.StoreID, url)
		}
	}

	wg.Wait()

	job.mu.Lock()
	if len(job.Errors) > 0 {
		job.Status = "failed"
	} else {
		job.Status = "completed"
	}
	job.mu.Unlock()
}

func processImage(url string) (ImageResult, error) {
	// Validate URL
	if url == "" {
		return ImageResult{}, fmt.Errorf("empty image URL")
	}

	// Download image
	resp, err := http.Get(url)
	if err != nil {
		return ImageResult{}, fmt.Errorf("failed to download image: %v", err)
	}
	defer resp.Body.Close()

	// Decode image configuration
	img, _, err := image.DecodeConfig(resp.Body)
	if err != nil {
		return ImageResult{}, fmt.Errorf("failed to decode image: %v", err)
	}

	// Calculate perimeter
	perimeter := 2 * (img.Width + img.Height)

	log.Printf("Processing image from %s, Dimensions: %dx%d, Perimeter: %d", url, img.Width, img.Height, perimeter)
	// Simulate GPU processing
	sleepTime := time.Duration(100+rand.Intn(300)) * time.Millisecond
	time.Sleep(sleepTime)

	// Create and return image result
	return ImageResult{
		URL:       url,
		Width:     img.Width,
		Height:    img.Height,
		Perimeter: perimeter,
	}, nil
}
