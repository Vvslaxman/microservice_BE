package main

import (
	"fmt"
	"log"
	"net/http"
	"kiranaC_BE/models"
	"kiranaC_BE/routes"
)

func main() {
	var err error
	// Load store data into routes.StoreData
	routes.StoreData, err = models.LoadStores("D://kiranaC_BE//data//StoreMasterAssignment.csv")
	if err != nil {
		log.Fatalf("Failed to load store data: %v", err)
	}

	http.HandleFunc("/api/submit", routes.SubmitJobHandler)
	http.HandleFunc("/api/status", routes.StatusHandler)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
