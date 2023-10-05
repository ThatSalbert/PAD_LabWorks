package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"weather-service/database"
	"weather-service/payload"
)

var db *sql.DB
var err error

var maxConnections = 10

func registerService(HOSTNAME string, PORT string) {
	var jsonRequest = []byte(`{
		"service_name": "weather",
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

	//GET /weather/locations?country={country}
	router.HandleFunc("/weather/locations", GetLocations).Methods("GET").Queries("country", "{country}")

	//GET /weather/current?city={city}
	router.HandleFunc("/weather/current", GetCurrentWeather).Methods("GET").Queries("city", "{city}")

	//GET /weather/forecast?city={city}
	router.HandleFunc("/weather/forecast", GetWeatherForecast).Methods("GET").Queries("city", "{city}")

	//POST /weather/add_data?type={type}
	router.HandleFunc("/weather/add_data", AddWeatherData).Methods("POST").Queries("type", "{type}")

	//DELETE /weather/delete_data?type={type}
	router.HandleFunc("/weather/delete_data", DeleteWeatherData).Methods("DELETE").Queries("type", "{type}")

	registerService("localhost", "8000")

	log.Fatal(http.ListenAndServe(":8000", router))

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
	}
	locations, funcErr := database.GetCities(country, db)
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
	city := query.Get("city")
	if len(city) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"message": "location query parameter not specified"}`))
		if err != nil {
			return
		}
	}
	weather, funcErr := database.GetCurrentWeather(city, db)
	if funcErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": "Internal server error"}`))
		if err != nil {
			return
		}
	} else {
		if len(weather) == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte(`{"message": "No weather data found for location = ` + city + `"}`))
			if err != nil {
				return
			}
		} else {
			response := payload.GenerateCurrentWeatherResponse(weather[0].Country, weather[0].Location, weather[0].Timestamp, weather[0].Temperature, weather[0].Humidity, weather[0].WeatherCondition)
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

func GetWeatherForecast(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	city := query.Get("city")
	if len(city) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"message": "location query parameter not specified"}`))
		if err != nil {
			return
		}
	}
	forecast, funcErr := database.GetForecastWeather(city, db)
	if funcErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": "Internal server error"}`))
		if err != nil {
			return
		}
	} else {
		if len(forecast) == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte(`{"message": "No weather data found for location = ` + city + `"}`))
			if err != nil {
				return
			}
		} else {
			response := payload.GenerateForecastWeatherResponse(forecast[0].Country, forecast[0].Location, forecast[0].ForecastWeatherDays)
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
			_, err := w.Write([]byte(`{"message": "Invalid JSON payload"}`))
			if err != nil {
				return
			}
		}
		funcErr := database.AddCurrentWeather(weatherData, db)
		if funcErr != nil {
			if funcErr.Error() == "value already exists" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusConflict)
				_, err := w.Write([]byte(`{"message": "Value already exists"}`))
				if err != nil {
					return
				}
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write([]byte(`{"message": "Internal server error"}`))
				if err != nil {
					return
				}
			}
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"message": "AddWeatherData called with type=` + dataType + `"}`))
			if err != nil {
				return
			}
		}
	} else if dataType == "forecast" {
		var forecastData payload.AddDataForecast
		err := jsonDecoder.Decode(&forecastData)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(`{"message": "Invalid JSON payload"}`))
			if err != nil {
				return
			}
		}
		funcErr := database.AddForecastWeather(forecastData, db)
		if funcErr != nil {
			if funcErr.Error() == "value already exists" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusConflict)
				_, err := w.Write([]byte(`{"message": "Value already exists"}`))
				if err != nil {
					return
				}
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write([]byte(`{"message": "Internal server error"}`))
				if err != nil {
					return
				}
			}
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"message": "AddWeatherData called with type=` + dataType + `"}`))
			if err != nil {
				return
			}
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"message": "Invalid type query parameter"}`))
		if err != nil {
			return
		}
	}
}

func DeleteWeatherData(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	dataType := query.Get("type")
	jsonDecoder := json.NewDecoder(r.Body)
	if dataType == "weather" {
		var weatherData payload.DeleteData
		err := jsonDecoder.Decode(&weatherData)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(`{"message": "Invalid JSON payload"}`))
			if err != nil {
				return
			}
		}
		funcErr := database.DeleteWeather(weatherData.City, weatherData.Timestamp, db)
		if funcErr != nil {
			if funcErr.Error() == "value does not exist" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				_, err := w.Write([]byte(`{"message": "Value does not exist"}`))
				if err != nil {
					return
				}
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write([]byte(`{"message": "Internal server error"}`))
				if err != nil {
					return
				}
			}
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"message": "DeleteWeatherData called with type=` + dataType + `"}`))
			if err != nil {
				return
			}
		}
	} else if dataType == "forecast" {
		var forecastData payload.DeleteData
		err := jsonDecoder.Decode(&forecastData)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(`{"message": "Invalid JSON payload"}`))
			if err != nil {
				return
			}
		}
		funcErr := database.DeleteForecast(forecastData.City, forecastData.Timestamp, db)
		if funcErr != nil {
			if funcErr.Error() == "value does not exist" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				_, err := w.Write([]byte(`{"message": "Value does not exist"}`))
				if err != nil {
					return
				}
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write([]byte(`{"message": "Internal server error"}`))
				if err != nil {
					return
				}
			}
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"message": "DeleteWeatherData called with type=` + dataType + `"}`))
			if err != nil {
				return
			}
		}
	}
}
