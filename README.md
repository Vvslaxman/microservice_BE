# Image Processing API

## Overview
This project provides two endpoints for image processing:
1. **Submit Job** - Submit images for processing.
2. **Get Job Info** - Retrieve the status of a submitted job.

## Directory Structure
.
├── main.go
├── routes/
│   ├── submit.go
│   ├── status.go
├── models/
│   ├── job.go
│   ├── store.go
├── utils/
│   ├── downloader.go
│   ├── worker.go
├── data/
│   ├── StoreMasterAssignment.csv
├── Dockerfile
├── docker-compose.yml
├── README.md


## Usage
### Build and Run with Docker
1. Build the image:
   ```sh
   docker-compose build

2. Run the application:
   ```sh
   docker-compose up

### Endpoints
Submit Job: POST /api/submit/
Get Job Info: GET /api/status?jobid={id}

