package main

import (
	"bytes"
	"database/sql"
	"disaster-service/database"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var db *sql.DB
var err error

var maxConnections = 10

func registerService(HOSTNAME string, PORT string, SERVICE_NAME string) {
	var jsonRequest = []byte(`{
		"service_name": "` + SERVICE_NAME + `",
		"service_address": "http://` + HOSTNAME + `:` + PORT + `"
	}`)
	response, err := http.Post("http://localhost:8002/register", "application/json", bytes.NewBuffer(jsonRequest))
	if err != nil {
		log.Fatal(err)
	} else {
		if response.StatusCode == 200 {
			log.Println("Service registered successfully")
		} else if response.StatusCode == 409 {
			log.Println("Service already registered")
		} else {
			log.Println("Service registration failed")
		}
	}
}

func main() {
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

	//GET /disaster
	router.HandleFunc("/disaster", GetDisasters).Methods("GET")

	//GET /disaster/list?city={city}&active={active}
	router.HandleFunc("/disaster/list", GetDisasterList).Methods("GET").Queries("city", "{city}", "active", "{active}")

	//POST /disaster/alert?city={city}
	router.HandleFunc("/disaster/alert", PostDisasterAlert).Methods("POST").Queries("city", "{city}")

	//PUT /disaster/alert?alert_id={alert_id}
	router.HandleFunc("/disaster/alert", PutDisasterAlert).Methods("PUT").Queries("alert_id", "{alert_id}")

	//DELETE /disaster/alert?alert_id={alert_id}
	router.HandleFunc("/disaster/alert", DeleteDisasterAlert).Methods("DELETE").Queries("alert_id", "{alert_id}")

	registerService("localhost", "8001", "disaster")

	log.Fatal(http.ListenAndServe(":8001", router))
}

func GetDisasters(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"message": "GetDisasters called"}`))
	if err != nil {
		return
	}
}

func GetDisasterList(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	location := query.Get("city")
	active := query.Get("active")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"message": "GetDisasterList called with location=` + location + ` and active=` + active + `"}`))
	if err != nil {
		return
	}
}

func PostDisasterAlert(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	location := query.Get("city")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"message": "PostDisasterAlert called with location=` + location + `"}`))
	if err != nil {
		return
	}
}

func PutDisasterAlert(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	alertId := query.Get("alert_id")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"message": "PutDisasterAlert called with alert_id=` + alertId + `"}`))
	if err != nil {
		return
	}
}

func DeleteDisasterAlert(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	alertId := query.Get("alert_id")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"message": "DeleteDisasterAlert called with alert_id=` + alertId + `"}`))
	if err != nil {
		return
	}
}
