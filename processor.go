package main

import (
    "image"
    _ "image/jpeg"
    _ "image/png"
    "math/rand"
    "net/http"
    "time"
	"sync"
	"fmt"
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

                if err := processImage(imageURL); err != nil {
                    job.mu.Lock()
                    job.Errors = append(job.Errors, ErrorDetail{
                        StoreID: storeID,
                        Error:   err.Error(),
                    })
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

func processImage(url string) error {
    // Download and get image dimensions
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    img, _, err := image.DecodeConfig(resp.Body)
    if err != nil {
        return err
    }

    // Calculate perimeter
    perimeter := 2 * (img.Width + img.Height)

    // Use perimeter (e.g., log it)
    fmt.Printf("Processing image from %s, Perimeter: %d\n", url, perimeter)

    // Simulate GPU processing
    time.Sleep(time.Duration(100+rand.Intn(300)) * time.Millisecond)

    return nil
}
