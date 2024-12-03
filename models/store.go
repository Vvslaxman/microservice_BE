package models

import (
	"encoding/csv"
	"os"
	"errors"
)

type Store struct {
	StoreID string
	Name    string
	Address string
}

var StoreData map[string]string
// LoadStoreData loads store data from the CSV file into a map for validation
func LoadStores(filePath string) (map[string]string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    rows, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }

    if len(rows) < 2 {
        return nil, errors.New("CSV file is empty or missing header")
    }

    // Dynamically identify column indices from header row
    header := rows[0]
    storeIDIndex, storeNameIndex := -1, -1
    for i, col := range header {
        if col == "StoreID" {
            storeIDIndex = i
        } else if col == "StoreName" {
            storeNameIndex = i
        }
    }

    if storeIDIndex == -1 || storeNameIndex == -1 {
        return nil, errors.New("StoreID or StoreName column missing in CSV")
    }

    // Load data into map
    storeData := make(map[string]string)
    for _, row := range rows[1:] {
        storeID := row[storeIDIndex]
        storeName := row[storeNameIndex]
        storeData[storeID] = storeName
    }

    

    return storeData, nil
}