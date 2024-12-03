package main

import (
    "encoding/csv"
    "os"
	"fmt"
	"log"
)

var storeData = make(map[string]bool)

func LoadStoreData(filepath string) error {
    file, err := os.Open(filepath)
    if err != nil {
        return err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        return err
    }

	// Find StoreID column index
    var storeIDIndex int = -1
    headers := records[0]
    for i, header := range headers {
        if header == "StoreID" {
            storeIDIndex = i
            break
        }
    }

    if storeIDIndex == -1 {
        return fmt.Errorf("StoreID column not found in CSV")
    }

    for i, record := range records {
        if i == 0 { // Skip header
            continue
        }
        storeID := record[storeIDIndex]
        storeData[storeID] = true
        log.Printf("Loaded store ID: %s", storeID) // Debug log
    }
	log.Printf("Loaded %d stores", len(storeData))
    return nil
}

func ValidateStore(storeID string) bool {
    valid := storeData[storeID]
    log.Printf("Validating store ID %s: %v", storeID, valid) // Debug log
    return valid
}