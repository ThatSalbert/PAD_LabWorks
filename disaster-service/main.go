package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	//GET /disasters
	router.HandleFunc("/disasters", GetDisasters).Methods("GET")

	//GET /disasters/list?location={location}&active={active}
	router.HandleFunc("/disasters/list", GetDisasterList).Methods("GET").Queries("location", "{location}", "active", "{active}")

	//POST /disasters/alert?location={location}
	router.HandleFunc("/disasters/alert", PostDisasterAlert).Methods("POST").Queries("location", "{location}")

	//PUT /disasters/alert/{alert_id}
	router.HandleFunc("/disasters/alert/{alert_id}", PutDisasterAlert).Methods("PUT")

	//DELETE /disasters/alert/{alert_id}
	router.HandleFunc("/disasters/alert/{alert_id}", DeleteDisasterAlert).Methods("DELETE")

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
	location := query.Get("location")
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
	location := query.Get("location")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"message": "PostDisasterAlert called with location=` + location + `"}`))
	if err != nil {
		return
	}
}

func PutDisasterAlert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	alertId := vars["alert_id"]
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"message": "PutDisasterAlert called with alert_id=` + alertId + `"}`))
	if err != nil {
		return
	}
}

func DeleteDisasterAlert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	alertId := vars["alert_id"]
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"message": "DeleteDisasterAlert called with alert_id=` + alertId + `"}`))
	if err != nil {
		return
	}
}
