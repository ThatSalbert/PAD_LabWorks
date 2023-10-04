package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"weather-service/database"
	"weather-service/payload"
	"weather-service/worker"
)

var db *sql.DB
var err error

var maxConnections = 10

func main() {
	worker.InitWorkers(maxConnections)

	db, err = database.ConnectToDB(maxConnections)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	router := mux.NewRouter()

	//GET /weather/locations?country={country}
	router.HandleFunc("/weather/locations", GetLocations).Methods("GET").Queries("country", "{country}")

	//GET /weather/current?location={location}
	router.HandleFunc("/weather/current", GetCurrentWeather).Methods("GET").Queries("location", "{location}")

	//GET /weather/forecast?location={location}
	router.HandleFunc("/weather/forecast", GetWeatherForecast).Methods("GET").Queries("location", "{location}")

	//POST /weather/add_data?type={type}
	router.HandleFunc("/weather/add_data", AddWeatherData).Methods("POST").Queries("type", "{type}")

	//DELETE /weather/delete_data?type={type}
	router.HandleFunc("/weather/delete_data", DeleteWeatherData).Methods("DELETE").Queries("type", "{type}")

	log.Fatal(http.ListenAndServe(":8000", router))
}
func GetLocations(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	country := query.Get("country")

	responseChannel := make(chan worker.LocationResponse)
	worker.LocationRequestQueue <- worker.LocationRequest{Country: country, DB: db, Response: responseChannel}

	response := <-responseChannel
	locations := response.Result
	funcErr := response.Err

	if funcErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": "Internal server error"}`))
		if err != nil {
			return
		}
	} else {
		if len(locations) == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte(`{"message": "No locations found for country = ` + country + `"}`))
			if err != nil {
				return
			}
		} else {
			response := payload.GenerateLocationResponse(country, locations)
			jsonResponse, _ := json.Marshal(response)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write(jsonResponse)
			if err != nil {
				return
			}
		}
	}
}

func GetCurrentWeather(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	location := query.Get("location")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"message": "GetCurrentWeather called with location=` + location + `"}`))
	if err != nil {
		return
	}
}

func GetWeatherForecast(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	location := query.Get("location")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"message": "GetWeatherForecast called with location=` + location + `"}`))
	if err != nil {
		return
	}
}

func AddWeatherData(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	dataType := query.Get("type")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"message": "AddWeatherData called with type=` + dataType + `"}`))
	if err != nil {
		return
	}
}

func DeleteWeatherData(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	dataType := query.Get("type")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"message": "DeleteWeatherData called with type=` + dataType + `"}`))
	if err != nil {
		return
	}
}
