package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"weather-service/database"
	"weather-service/payload"

	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

var maxConnections = 10

var (
	WEATHER_HOSTNAME     = os.Getenv("WEATHER_HOSTNAME")
	WEATHER_PORT         = os.Getenv("WEATHER_PORT")
	SERVICEDISC_HOSTNAME = os.Getenv("SERVICEDISC_HOSTNAME")
	SERVICEDISC_PORT     = os.Getenv("SERVICEDISC_PORT")
	DB_HOST              = os.Getenv("DB_HOST")
	DB_PORT              = os.Getenv("DB_PORT")
)

func registerService(WEATHER_HOSTNAME string, WEATHER_PORT string, SERVICEDISC_HOSTNAME string, SERVICEDISC_PORT string) {
	var jsonRequest = []byte(`{
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

	registerService(WEATHER_HOSTNAME, WEATHER_PORT, SERVICEDISC_HOSTNAME, SERVICEDISC_PORT)

	log.Fatal(http.ListenAndServe(":"+WEATHER_PORT, router))

}
func GetLocations(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	country := query.Get("country")
	if country == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"message": "country query parameter not specified"}`))
		if err != nil {
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

func GetCurrentWeather(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	country := query.Get("country")
	city := query.Get("city")
	if len(country) == 0 || len(city) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"message": "country or city query parameter not specified"}`))
		if err != nil {
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
		if err != nil {
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
			if err != nil {
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
		if err != nil {
			return
		}
		return
	case 404:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		var funcErrString string
		funcErrString = funcErr.Error()
		_, err := w.Write([]byte(`{"message": ` + "\"" + funcErrString + "\"" + `}`))
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

func GetWeatherForecast(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	country := query.Get("country")
	city := query.Get("city")
	if len(country) == 0 || len(city) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"message": "country or city query parameter not specified"}`))
		if err != nil {
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

func AddWeatherData(w http.ResponseWriter, r *http.Request) {
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
			if err != nil {
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
			if err != nil {
				return
			}
			return
		case 409:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			_, err := w.Write([]byte(`{"message": ` + "\"" + funcErr.Error() + "\"" + `}`))
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
	} else if dataType == "forecast" {
		var forecastData payload.AddDataForecast
		err := jsonDecoder.Decode(&forecastData)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(`{"message": "invalid JSON payload"}`))
			if err != nil {
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
			if err != nil {
				return
			}
			return
		case 409:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			_, err := w.Write([]byte(`{"message": ` + "\"" + funcErr.Error() + "\"" + `}`))
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
}

func UpdateWeatherData(w http.ResponseWriter, r *http.Request) {
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
			if err != nil {
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
			if err != nil {
				return
			}
			return
		case 404:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte(`{"message": "weather data not found"}`))
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
	} else if dataType == "forecast" {
		var forecastData payload.UpdateDataForecast
		err := jsonDecoder.Decode(&forecastData)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(`{"message": "invalid JSON payload"}`))
			if err != nil {
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
			if err != nil {
				return
			}
			return
		case 404:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte(`{"message": "forecast data not found"}`))
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
}
