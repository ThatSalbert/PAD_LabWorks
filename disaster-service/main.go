package main

import (
	"bytes"
	"database/sql"
	"disaster-service/database"
	"disaster-service/payload"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	_ "github.com/prometheus/client_golang/prometheus"
	_ "github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

var maxConnections = 10

var (
	SERVICE_TYPE         = os.Getenv("SERVICE_TYPE")
	DISASTER_HOSTNAME    = os.Getenv("DISASTER_HOSTNAME")
	DISASTER_PORT        = os.Getenv("DISASTER_PORT")
	SERVICEDISC_HOSTNAME = os.Getenv("SERVICEDISC_HOSTNAME")
	SERVICEDISC_PORT     = os.Getenv("SERVICEDISC_PORT")
	DB_HOST              = os.Getenv("DB_HOST")
	DB_PORT              = os.Getenv("DB_PORT")
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"code", "method"},
	)

	httpRequestsDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_requests_duration_seconds",
			Help: "Duration of HTTP requests",
		},
		[]string{"code", "method"},
	)

	errCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "error_counter",
			Help: "Total number of errors",
		},
		[]string{"code", "method"},
	)
)

func registerService(DISASTER_HOSTNAME string, DISASTER_PORT string, SERVICEDISC_HOSTNAME string, SERVICEDISC_PORT string, SERVICE_TYPE string) {
	var jsonRequest = []byte(`{
		"service_type": "` + SERVICE_TYPE + `",
		"service_name": "` + DISASTER_HOSTNAME + `",
		"service_address": "http://` + DISASTER_HOSTNAME + `:` + DISASTER_PORT + `"
	}`)
	response, err := http.Post("http://"+SERVICEDISC_HOSTNAME+":"+SERVICEDISC_PORT+"/register", "application/json", bytes.NewBuffer(jsonRequest))
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
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestsDuration)
	prometheus.MustRegister(errCounter)

	db, err = database.ConnectToDB(maxConnections, DB_HOST, DB_PORT)
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

	registerService(DISASTER_HOSTNAME, DISASTER_PORT, SERVICEDISC_HOSTNAME, SERVICEDISC_PORT, SERVICE_TYPE)

	log.Fatal(http.ListenAndServe(":"+DISASTER_PORT, router))
}

func GetDisasters(w http.ResponseWriter, r *http.Request) {
	httpRequestsTotal.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Inc()
	startTimer := time.Now()
	disasters, funcErrCode, funcErr := database.GetDisasterTypes(db)
	switch funcErrCode {
	case 200:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(disasters)
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Inc()
			return
		}
		return
	case 404:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"message": ` + "\"" + funcErr.Error() + "\"" + `}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusNotFound), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Inc()
			return
		}
		return
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": ` + "\"" + funcErr.Error() + "\"" + `}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Inc()
			return
		}
		return
	}
}

func GetDisasterList(w http.ResponseWriter, r *http.Request) {
	httpRequestsTotal.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Inc()
	startTimer := time.Now()
	query := r.URL.Query()
	country := query.Get("country")
	city := query.Get("city")
	active := query.Get("active")
	if active == "" || city == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"message": "country query parameter not specified"}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Inc()
			return
		}
		return
	}
	var activeBool = false
	if active == "true" {
		activeBool = true
	} else if active == "false" {
		activeBool = false
	}
	disasterList, funcErrCode, funcErr := database.GetDisasterList(db, country, city, activeBool)
	switch funcErrCode {
	case 200:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(disasterList)
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Inc()
			return
		}
		return
	case 404:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"message": ` + "\"" + funcErr.Error() + "\"" + `}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusNotFound), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Inc()
			return
		}
		return
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": ` + "\"" + funcErr.Error() + "\"" + `}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Inc()
			return
		}
		return
	}
}

func PostDisasterAlert(w http.ResponseWriter, r *http.Request) {
	httpRequestsTotal.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Inc()
	startTimer := time.Now()
	jsonDecoder := json.NewDecoder(r.Body)
	var alert payload.AddAlert
	err := jsonDecoder.Decode(&alert)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"message": "invalid JSON payload"}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Inc()
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
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Inc()
			return
		}
		return
	case 409:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, err := w.Write([]byte(`{"message": "alert data already exists"}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusConflict), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Inc()
			return
		}
		return
	case 404:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"message": ` + "\"" + funcErr.Error() + "\"" + `}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusNotFound), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Inc()
			return
		}
		return
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": ` + "\"" + funcErr.Error() + "\"" + `}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Inc()
			return
		}
		return
	}
}

func PutDisasterAlert(w http.ResponseWriter, r *http.Request) {
	httpRequestsTotal.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Inc()
	startTimer := time.Now()
	query := r.URL.Query()
	alertId := query.Get("alert_id")
	if alertId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"message": "alert_id query parameter not specified"}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Inc()
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
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Inc()
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
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Inc()
			return
		}
		return
	case 404:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"message": ` + "\"" + funcErr.Error() + "\"" + `}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusNotFound), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Inc()
			return
		}
		return
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": ` + "\"" + funcErr.Error() + "\"" + `}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Inc()
			return
		}
		return
	}
}
