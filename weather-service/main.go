package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"weather-service/database"
	"weather-service/payload"

	"github.com/gorilla/mux"
)

var db *sql.DB
var err error
var (
	SERVICE_TYPE         = os.Getenv("SERVICE_TYPE")
	WEATHER_HOSTNAME     = os.Getenv("WEATHER_HOSTNAME")
	WEATHER_PORT         = os.Getenv("WEATHER_PORT")
	SERVICEDISC_HOSTNAME = os.Getenv("SERVICEDISC_HOSTNAME")
	SERVICEDISC_PORT     = os.Getenv("SERVICEDISC_PORT")
	DB_HOST              = os.Getenv("DB_HOST")
	DB_PORT              = os.Getenv("DB_PORT")
	MAX_CONNECTIONS      = os.Getenv("MAX_CONNECTIONS")
	METRICS_PORT		 = os.Getenv("METRICS_PORT")
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

func registerService(WEATHER_HOSTNAME string, WEATHER_PORT string, SERVICEDISC_HOSTNAME string, SERVICEDISC_PORT string, SERVICE_TYPE string) {
	var jsonRequest = []byte(`{
		"service_type": "` + SERVICE_TYPE + `",
		"service_name": "` + WEATHER_HOSTNAME + `",
		"service_address": "http://` + WEATHER_HOSTNAME + `:` + WEATHER_PORT + `"
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

	maxConnections, _ := strconv.Atoi(MAX_CONNECTIONS)
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

	//GET /weather/locations?country={country}
	router.HandleFunc("/weather/locations", GetLocations).Methods("GET").Queries("country", "{country}")

	//GET /weather/current?country={country}&city={city}
	router.HandleFunc("/weather/current", GetCurrentWeather).Methods("GET").Queries("country", "{country}", "city", "{city}")

	//GET /weather/forecast?country={country}&city={city}
	router.HandleFunc("/weather/forecast", GetWeatherForecast).Methods("GET").Queries("country", "{country}", "city", "{city}")

	//POST /weather/add_data?type={type}
	router.HandleFunc("/weather/add_data", AddWeatherData).Methods("POST").Queries("type", "{type}")

	//PUT /weather/update_data?type={type}
	router.HandleFunc("/weather/update_data", UpdateWeatherData).Methods("PUT").Queries("type", "{type}")

	//POST /weather/add_city/prepare
	router.HandleFunc("/weather/add_city/prepare", AddCityPrepare).Methods("POST")

	//POST /weather/add_city/commit
	router.HandleFunc("/weather/add_city/commit", AddCity).Methods("POST")

	//POST /weather/add_city/rollback
	router.HandleFunc("/weather/add_city/rollback", AddCityRollback).Methods("POST")

	router.Handle("/metrics", promhttp.Handler())

	go func() {
		if err := http.ListenAndServe(":"+METRICS_PORT, nil); err != nil {
			log.Fatal(err)
		}
	}()

	registerService(WEATHER_HOSTNAME, WEATHER_PORT, SERVICEDISC_HOSTNAME, SERVICEDISC_PORT, SERVICE_TYPE)

	log.Fatal(http.ListenAndServe(":"+WEATHER_PORT, router))

}
func GetLocations(w http.ResponseWriter, r *http.Request) {
	httpRequestsTotal.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Inc()
	startTimer := time.Now()
	query := r.URL.Query()
	country := query.Get("country")
	if country == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"message": "country query parameter not specified"}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Inc()
			return
		}
		return
	}
	locations, funcErrCode, funcErr := database.GetCities(country, db)
	switch funcErrCode {
	case 200:
		var locationResponse payload.LocationResponse
		locationResponse.Country = country
		locationResponse.Locations = locations
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		jsonResponse, _ := json.Marshal(locationResponse)
		_, err := w.Write(jsonResponse)
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Inc()
			return
		}
		return
	case 404:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"message": ` + "\"" + funcErr.Error() + "\"" + `}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusNotFound), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusNotFound), r.Method).Inc()
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

func GetCurrentWeather(w http.ResponseWriter, r *http.Request) {
	httpRequestsTotal.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Inc()
	startTimer := time.Now()
	query := r.URL.Query()
	country := query.Get("country")
	city := query.Get("city")
	if len(country) == 0 || len(city) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"message": "country or city query parameter not specified"}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Inc()
			return
		}
		return
	}

	disasterAddress := "http://disaster-service-1:9091/disaster/list"
	list, err := http.Get(disasterAddress + "?country=" + country + "&city=" + city + "&active=true")
	if err != nil {
		fmt.Println("Error getting disaster list")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": "internal server error"}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Inc()
			return
		}
		return
	}
	var disasters []payload.Disaster
	if list.StatusCode == 404 {
		fmt.Println("Disaster list not found")
		disasters = nil
	} else if list.StatusCode == 200 {
		jsonDecoder := json.NewDecoder(list.Body)
		err := jsonDecoder.Decode(&disasters)
		if err != nil {
			fmt.Println("Error decoding disaster list")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(`{"message": "internal server error"}`))
			httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Observe(time.Since(startTimer).Seconds())
			if err != nil {
				errCounter.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method).Inc()
				return
			}
			return
		}
	}
	weather, funcCodeErr, funcErr := database.GetCurrentWeather(country, city, disasters, db)
	switch funcCodeErr {
	case 200:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		jsonResponse, _ := json.Marshal(weather)
		_, err := w.Write(jsonResponse)
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Inc()
			return
		}
		return
	case 404:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		var funcErrString string
		funcErrString = funcErr.Error()
		_, err := w.Write([]byte(`{"message": ` + "\"" + funcErrString + "\"" + `}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusNotFound), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusNotFound), r.Method).Inc()
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

func GetWeatherForecast(w http.ResponseWriter, r *http.Request) {
	httpRequestsTotal.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Inc()
	startTimer := time.Now()
	query := r.URL.Query()
	country := query.Get("country")
	city := query.Get("city")
	if len(country) == 0 || len(city) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"message": "country or city query parameter not specified"}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Inc()
			return
		}
		return
	}
	weather, funcCodeErr, funcErr := database.GetForecastWeather(country, city, db)
	switch funcCodeErr {
	case 200:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		jsonResponse, _ := json.Marshal(weather)
		_, err := w.Write(jsonResponse)
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Inc()
			return
		}
		return
	case 404:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"message": ` + "\"" + funcErr.Error() + "\"" + `}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusNotFound), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusNotFound), r.Method).Inc()
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

func AddWeatherData(w http.ResponseWriter, r *http.Request) {
	httpRequestsTotal.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Inc()
	startTimer := time.Now()
	query := r.URL.Query()
	dataType := query.Get("type")
	jsonDecoder := json.NewDecoder(r.Body)
	if dataType == "weather" {
		var weatherData payload.AddDataWeather
		err := jsonDecoder.Decode(&weatherData)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(`{"message": "invalid JSON payload"}`))
			httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Observe(time.Since(startTimer).Seconds())
			if err != nil {
				errCounter.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Inc()
				return
			}
			return
		}
		weatherID, funcErrCode, funcErr := database.AddCurrentWeather(weatherData, db)
		switch funcErrCode {
		case 200:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"message": "weather data added successfully with id=` + strconv.FormatInt(*weatherID, 10) + `"}`))
			httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Observe(time.Since(startTimer).Seconds())
			if err != nil {
				errCounter.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Inc()
				return
			}
			return
		case 409:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			_, err := w.Write([]byte(`{"message": ` + "\"" + funcErr.Error() + "\"" + `}`))
			httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusConflict), r.Method).Observe(time.Since(startTimer).Seconds())
			if err != nil {
				errCounter.WithLabelValues(strconv.Itoa(http.StatusConflict), r.Method).Inc()
				return
			}
			return
		case 404:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte(`{"message": ` + "\"" + funcErr.Error() + "\"" + `}`))
			httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusNotFound), r.Method).Observe(time.Since(startTimer).Seconds())
			if err != nil {
				errCounter.WithLabelValues(strconv.Itoa(http.StatusNotFound), r.Method).Inc()
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
	} else if dataType == "forecast" {
		var forecastData payload.AddDataForecast
		err := jsonDecoder.Decode(&forecastData)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(`{"message": "invalid JSON payload"}`))
			httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Observe(time.Since(startTimer).Seconds())
			if err != nil {
				errCounter.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Inc()
				return
			}
			return
		}
		forecastID, funcErrCode, funcErr := database.AddForecastWeather(forecastData, db)
		switch funcErrCode {
		case 200:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"message": "forecast data added successfully with id=` + strconv.FormatInt(*forecastID, 10) + `"}`))
			httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Observe(time.Since(startTimer).Seconds())
			if err != nil {
				errCounter.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Inc()
				return
			}
			return
		case 409:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			_, err := w.Write([]byte(`{"message": ` + "\"" + funcErr.Error() + "\"" + `}`))
			httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusConflict), r.Method).Observe(time.Since(startTimer).Seconds())
			if err != nil {
				errCounter.WithLabelValues(strconv.Itoa(http.StatusConflict), r.Method).Inc()
				return
			}
			return
		case 404:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte(`{"message": ` + "\"" + funcErr.Error() + "\"" + `}`))
			httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusNotFound), r.Method).Observe(time.Since(startTimer).Seconds())
			if err != nil {
				errCounter.WithLabelValues(strconv.Itoa(http.StatusNotFound), r.Method).Inc()
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
}

func UpdateWeatherData(w http.ResponseWriter, r *http.Request) {
	httpRequestsTotal.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Inc()
	startTimer := time.Now()
	query := r.URL.Query()
	dataType := query.Get("type")
	jsonDecoder := json.NewDecoder(r.Body)
	if dataType == "weather" {
		var weatherData payload.UpdateDataWeather
		err := jsonDecoder.Decode(&weatherData)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(`{"message": "invalid JSON payload"}`))
			httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Observe(time.Since(startTimer).Seconds())
			if err != nil {
				errCounter.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Inc()
				return
			}
			return
		}
		funcErrCode, funcErr := database.UpdateWeather(weatherData.WeatherID, weatherData, db)
		switch funcErrCode {
		case 200:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"message": "weather data updated successfully with id=` + strconv.FormatInt(weatherData.WeatherID, 10) + `"}`))
			httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Observe(time.Since(startTimer).Seconds())
			if err != nil {
				errCounter.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Inc()
				return
			}
			return
		case 404:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte(`{"message": "weather data not found"}`))
			httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusNotFound), r.Method).Observe(time.Since(startTimer).Seconds())
			if err != nil {
				errCounter.WithLabelValues(strconv.Itoa(http.StatusNotFound), r.Method).Inc()
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
	} else if dataType == "forecast" {
		var forecastData payload.UpdateDataForecast
		err := jsonDecoder.Decode(&forecastData)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(`{"message": "invalid JSON payload"}`))
			httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Observe(time.Since(startTimer).Seconds())
			if err != nil {
				errCounter.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Inc()
				return
			}
			return
		}
		funcErrCode, funcErr := database.UpdateForecast(forecastData.ForecastID, forecastData, db)
		switch funcErrCode {
		case 200:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"message": "forecast data updated successfully with id=` + strconv.FormatInt(forecastData.ForecastID, 10) + `"}`))
			httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Observe(time.Since(startTimer).Seconds())
			if err != nil {
				errCounter.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Inc()
				return
			}
			return
		case 404:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte(`{"message": "forecast data not found"}`))
			httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusNotFound), r.Method).Observe(time.Since(startTimer).Seconds())
			if err != nil {
				errCounter.WithLabelValues(strconv.Itoa(http.StatusNotFound), r.Method).Inc()
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
}

func AddCityPrepare(w http.ResponseWriter, r *http.Request) {
	httpRequestsTotal.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Inc()
	startTimer := time.Now()
	jsonDecoder := json.NewDecoder(r.Body)
	var cityData payload.AddCity
	err := jsonDecoder.Decode(&cityData)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"message": "invalid JSON payload"}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Inc()
			return
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"message": "city data prepared successfully", "payload": { "country": "` + cityData.Country + `", "city": "` + cityData.City + `", "latitude": ` + strconv.FormatFloat(float64(cityData.Latitude), 'f', -1, 64) + `, "longitude": ` + strconv.FormatFloat(float64(cityData.Longitude), 'f', -1, 64) + `}}`))
	httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Observe(time.Since(startTimer).Seconds())
	if err != nil {
		errCounter.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Inc()
		return
	}
	return
}

func AddCity(w http.ResponseWriter, r *http.Request) {
	httpRequestsTotal.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Inc()
	startTimer := time.Now()
	jsonDecoder := json.NewDecoder(r.Body)
	var cityData payload.AddCity
	err := jsonDecoder.Decode(&cityData)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"message": "invalid JSON payload"}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Inc()
			return
		}
		return
	}
	funcErrCode, funcErr := database.AddCity(cityData, db)
	switch funcErrCode {
	case 200:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"message": "city added successfully"}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Inc()
			return
		}
		return
	case 409:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, err := w.Write([]byte(`{"message": "city already exists"}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusConflict), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusConflict), r.Method).Inc()
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

func AddCityRollback(w http.ResponseWriter, r *http.Request) {
	httpRequestsTotal.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Inc()
	startTimer := time.Now()
	jsonDecoder := json.NewDecoder(r.Body)
	var cityData payload.RemoveCity
	err := jsonDecoder.Decode(&cityData)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"message": "invalid JSON payload"}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Inc()
			return
		}
		return
	}
	funcErrCode, funcErr := database.RemoveCity(cityData, db)
	switch funcErrCode {
	case 200:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"message": "city removed successfully"}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Inc()
			return
		}
		return
	case 404:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"message": "city not found"}`))
		httpRequestsDuration.WithLabelValues(strconv.Itoa(http.StatusNotFound), r.Method).Observe(time.Since(startTimer).Seconds())
		if err != nil {
			errCounter.WithLabelValues(strconv.Itoa(http.StatusNotFound), r.Method).Inc()
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
