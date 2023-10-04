package worker

import (
	"database/sql"
	"fmt"
	"sync"
	"weather-service/database"
	"weather-service/payload"
)

var wg sync.WaitGroup
var LocationRequestQueue = make(chan LocationRequest)

type LocationRequest struct {
	Country  string
	DB       *sql.DB
	Response chan<- LocationResponse
}

type LocationResponse struct {
	Result []payload.Location
	Err    error
}

func InitWorkers(maxWorkers int) {
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go StartWorker(i + 1)
	}
}

func StartWorker(ID int) {
	defer wg.Done()
	fmt.Print("Worker ", ID, " started\n")

	for {
		select {
		case job := <-LocationRequestQueue:
			fmt.Print("Worker ", ID, " received job\n")
			locations, funcErr := database.GetCities(job.Country, job.DB)
			if funcErr != nil {
				fmt.Print("Worker ", ID, ": GetCities called with error\n")
				job.Response <- LocationResponse{nil, funcErr}
			} else {
				fmt.Print("Worker ", ID, ": GetCities called successfully\n")
				job.Response <- LocationResponse{locations, nil}
			}
		}
	}
}
