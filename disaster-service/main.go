package main

import (
	"bytes"
	"database/sql"
	"disaster-service/database"
	"disaster-service/payload"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
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

	//GET /disaster/list?country={country}&city={city}&active={active}
	router.HandleFunc("/disaster/list", GetDisasterList).Methods("GET").Queries("country", "{country}", "city", "{city}", "active", "{active}")

	//POST /disaster/alert
	router.HandleFunc("/disaster/alert", PostDisasterAlert).Methods("POST")

	//PUT /disaster/alert?alert_id={alert_id}
	router.HandleFunc("/disaster/alert", PutDisasterAlert).Methods("PUT").Queries("alert_id", "{alert_id}")

	registerService("localhost", "8001", "disaster")

	log.Fatal(http.ListenAndServe(":8001", router))
}

func GetDisasters(w http.ResponseWriter, _ *http.Request) {
	disasters, funcErrCode, funcErr := database.GetDisasterTypes(db)
	switch funcErrCode {
	case 200:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(disasters)
		if err != nil {
			return
		}
		return
	case 404:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"message": ` + "\"" + funcErr.Error() + "\"" + `}`))
		if err != nil {
			return
		}
		return
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": ` + "\"" + funcErr.Error() + "\"" + `}`))
		if err != nil {
			return
		}
		return
	}
}

func GetDisasterList(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	country := query.Get("country")
	city := query.Get("city")
	active := query.Get("active")
	if active == "" || city == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"message": "country query parameter not specified"}`))
		if err != nil {
			return
		}
		return
	}
	var activeBool bool
	if active == "true" {
		activeBool = true
	} else if active == "false" {
		activeBool = false
	} else {
		activeBool = false
	}
	disasterList, funcErrCode, funcErr := database.GetDisasterList(db, country, city, activeBool)
	switch funcErrCode {
	case 200:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(disasterList)
		if err != nil {
			return
		}
		return
	case 404:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"message": ` + "\"" + funcErr.Error() + "\"" + `}`))
		if err != nil {
			return
		}
		return
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": ` + "\"" + funcErr.Error() + "\"" + `}`))
		if err != nil {
			return
		}
		return
	}
}

func PostDisasterAlert(w http.ResponseWriter, r *http.Request) {
	jsonDecoder := json.NewDecoder(r.Body)
	var alert payload.AddAlert
	err := jsonDecoder.Decode(&alert)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"message": "invalid JSON payload"}`))
		if err != nil {
			return
		}
		return
	}
	disasterTypeID, funcErrCode, funcErr := database.AddAlert(alert, db)
	switch funcErrCode {
	case 200:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"message": "alert data added successfully with alert_id=` + strconv.Itoa(*disasterTypeID) + `"}`))
		if err != nil {
			return
		}
		return
	case 409:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, err := w.Write([]byte(`{"message": "alert data already exists"}`))
		if err != nil {
			return
		}
		return
	case 404:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"message": ` + "\"" + funcErr.Error() + "\"" + `}`))
		if err != nil {
			return
		}
		return
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": ` + "\"" + funcErr.Error() + "\"" + `}`))
		if err != nil {
			return
		}
		return
	}
}

func PutDisasterAlert(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	alertId := query.Get("alert_id")
	if alertId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"message": "alert_id query parameter not specified"}`))
		if err != nil {
			return
		}
		return
	}
	jsonDecoder := json.NewDecoder(r.Body)
	var alert payload.UpdateAlert
	err := jsonDecoder.Decode(&alert)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"message": "invalid JSON payload"}`))
		if err != nil {
			return
		}
		return
	}
	alertID, _ := strconv.Atoi(alertId)
	funcErrCode, funcErr := database.UpdateAlert(alertID, alert, db)
	switch funcErrCode {
	case 200:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"message": "alert data updated successfully with alert_id=` + strconv.Itoa(alertID) + `"}`))
		if err != nil {
			return
		}
		return
	case 404:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"message": ` + "\"" + funcErr.Error() + "\"" + `}`))
		if err != nil {
			return
		}
		return
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": ` + "\"" + funcErr.Error() + "\"" + `}`))
		if err != nil {
			return
		}
		return
	}
}
