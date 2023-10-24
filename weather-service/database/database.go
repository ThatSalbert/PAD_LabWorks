package database

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"log"
	"time"
	"weather-service/payload"
)

func ConnectToDB(maxConnections int) (*sql.DB, error) {
	db, err := sql.Open("postgres", "user=postgres password=1234567890 dbname=weather-service-database sslmode=disable")
	if db == nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxConnections)
	db.SetMaxIdleConns(maxConnections / 2)
	db.SetConnMaxLifetime(time.Second * 5)
	return db, nil
}

func GetLocationID(country string, city string, db *sql.DB) (locationID *int, errCode int16, err error) {
	rows, err := db.Query("SELECT location_id FROM location_table WHERE UPPER(country) LIKE UPPER($1) AND UPPER(city) LIKE UPPER($2);", country, city)
	if err != nil {
		return nil, 500, errors.New("internal server error")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)
	if !rows.Next() {
		return nil, 404, errors.New("location not found")
	} else {
		if err = rows.Scan(&locationID); err != nil {
			return nil, 500, errors.New("internal server error")
		} else {
			return locationID, 200, nil
		}
	}
}

func GetCities(country string, db *sql.DB) (locations []payload.Location, errCode int16, err error) {
	rows, err := db.Query("SELECT city, longitude, latitude FROM location_table WHERE UPPER(country) LIKE UPPER($1)", country)
	if err != nil {
		return nil, 500, errors.New("internal server error")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)
	for rows.Next() {
		var location payload.Location
		if err = rows.Scan(&location.City, &location.Longitude, &location.Latitude); err != nil {
			return nil, 500, errors.New("internal server error")
		}
		if location.City != "" {
			locations = append(locations, location)
		} else {
			return nil, 404, errors.New("location not found")
		}
	}
	if len(locations) == 0 {
		return nil, 404, errors.New("location not found")
	} else {
		return locations, 200, nil
	}
}

func GetCurrentWeather(country string, city string, disasterList []payload.Disaster, db *sql.DB) (weather []payload.CurrentWeatherResponse, errCode int16, err error) {
	locationID, locationErrorCode, locationError := GetLocationID(country, city, db)
	if locationID == nil {
		return nil, locationErrorCode, locationError
	} else {
		rows, err := db.Query("SELECT lt.country, lt.city, lt.longitude, lt.latitude, cwt.timestamp, cwt.temperature, cwt.humidity, cwt.weather_condition FROM current_weather_table AS cwt INNER JOIN location_table AS lt ON lt.location_id = cwt.location_id WHERE cwt.location_id = $1;", locationID)
		if err != nil {
			return nil, 500, errors.New("internal server error")
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(rows)
		if errors.Is(rows.Scan(), sql.ErrNoRows) {
			return nil, 404, errors.New("no current weather found")
		}
		var currentWeather payload.CurrentWeatherResponse
		for rows.Next() {
			if err = rows.Scan(&currentWeather.Country, &currentWeather.Location.City, &currentWeather.Location.Longitude, &currentWeather.Location.Latitude, &currentWeather.Timestamp, &currentWeather.Temperature, &currentWeather.Humidity, &currentWeather.WeatherCondition); err != nil {
				return nil, 500, errors.New("internal server error")
			}
		}
		if currentWeather.Location.City == "" {
			return nil, 404, errors.New("no current weather found")
		} else {
			if len(disasterList) != 0 {
				if disasterList[0].DisasterName != "" {
					currentWeather.Disasters = disasterList
				} else {
					currentWeather.Disasters = nil
				}
			} else {
				currentWeather.Disasters = nil
			}
			weather = append(weather, currentWeather)
			return weather, 200, nil
		}
	}

}

func GetForecastWeather(country string, city string, db *sql.DB) (forecast []payload.ForecastWeatherResponse, errCode int16, err error) {
	locationID, locationErrorCode, locationError := GetLocationID(country, city, db)
	if locationID == nil {
		return nil, locationErrorCode, locationError
	} else {
		rows, err := db.Query("SELECT lt.country, lt.city, lt.longitude, lt.latitude, fwt.timestamp, fwt.temperature_high, fwt.temperature_low, fwt.humidity, fwt.weather_condition FROM forecast_weather_table AS fwt INNER JOIN location_table AS lt ON lt.location_id = fwt.location_id WHERE fwt.location_id = $1 AND fwt.timestamp >= NOW() AND fwt.timestamp <= NOW()::DATE + INTERVAL '5 days' ORDER BY lt.location_id, fwt.timestamp;", locationID)
		if err != nil {
			return nil, 500, errors.New("internal server error")
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(rows)
		if errors.Is(rows.Scan(), sql.ErrNoRows) {
			return nil, 404, errors.New("no forecast weather found")
		}
		var forecastWeather payload.ForecastWeatherResponse
		for rows.Next() {
			var forecastWeatherDay payload.ForecastWeatherDay
			if err = rows.Scan(&forecastWeather.Country, &forecastWeather.Location.City, &forecastWeather.Location.Longitude, &forecastWeather.Location.Latitude, &forecastWeatherDay.Timestamp, &forecastWeatherDay.TemperatureHigh, &forecastWeatherDay.TemperatureLow, &forecastWeatherDay.Humidity, &forecastWeatherDay.WeatherCondition); err != nil {
				return nil, 500, errors.New("internal server error")
			}
			forecastWeather.ForecastWeatherDays = append(forecastWeather.ForecastWeatherDays, forecastWeatherDay)
		}
		if len(forecastWeather.ForecastWeatherDays) == 0 {
			return nil, 404, errors.New("no forecast weather found")
		} else {
			forecast = append(forecast, forecastWeather)
			return forecast, 200, nil
		}
	}
}

func AddCurrentWeather(weather payload.AddDataWeather, db *sql.DB) (weatherID *int64, errCode int16, err error) {
	locationID, locationErrorCode, locationError := GetLocationID(weather.Country, weather.City, db)
	if locationID == nil {
		return nil, locationErrorCode, locationError
	} else {
		row, err := db.Query("INSERT INTO current_weather_table (location_id, timestamp, temperature, humidity, weather_condition) SELECT $1, $2, $3, $4, $5 WHERE NOT EXISTS (SELECT 1 FROM current_weather_table WHERE location_id = $1 AND timestamp = $2) RETURNING weather_id;", locationID, weather.Timestamp, weather.Temperature, weather.Humidity, weather.WeatherCondition)
		if err != nil {
			return nil, 500, errors.New("internal server error")
		}
		defer func(row *sql.Rows) {
			err := row.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(row)
		if !row.Next() {
			return nil, 409, errors.New("weather already exists")
		} else {
			if err = row.Scan(&weatherID); err != nil {
				return nil, 500, errors.New("internal server error")
			} else {
				return weatherID, 200, nil
			}
		}
	}
}

func AddForecastWeather(weather payload.AddDataForecast, db *sql.DB) (forecastID *int64, errCode int16, err error) {
	locationID, locationErrorCode, locationError := GetLocationID(weather.Country, weather.City, db)
	if locationID == nil {
		return nil, locationErrorCode, locationError
	} else {
		row, err := db.Query("INSERT INTO forecast_weather_table (location_id, timestamp, temperature_high, temperature_low, humidity, weather_condition) SELECT $1, $2, $3, $4, $5, $6 WHERE NOT EXISTS (SELECT 1 FROM forecast_weather_table WHERE location_id = $1 AND timestamp = $2) RETURNING forecast_id;", locationID, weather.Timestamp, weather.TemperatureHigh, weather.TemperatureLow, weather.Humidity, weather.WeatherCondition)
		if err != nil {
			return nil, 500, errors.New("internal server error")
		}
		defer func(row *sql.Rows) {
			err := row.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(row)
		if !row.Next() {
			return nil, 409, errors.New("weather already exists")
		} else {
			if err = row.Scan(&forecastID); err != nil {
				return nil, 500, errors.New("internal server error")
			} else {
				return forecastID, 200, nil
			}
		}
	}
}

func UpdateWeather(weatherID int64, weather payload.UpdateDataWeather, db *sql.DB) (errCode int16, err error) {
	locationID, locationErrorCode, locationError := GetLocationID(weather.Country, weather.City, db)
	if locationID == nil {
		return locationErrorCode, locationError
	} else {
		row, err := db.Query("UPDATE current_weather_table SET location_id = $1, timestamp = $2, temperature = $3, humidity = $4, weather_condition = $5 WHERE weather_id = $6 RETURNING *;", locationID, weather.Timestamp, weather.Temperature, weather.Humidity, weather.WeatherCondition, weatherID)
		if err != nil {
			return 500, errors.New("internal server error")
		}
		defer func(row *sql.Rows) {
			err := row.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(row)
		if !row.Next() {
			return 404, errors.New("weather not found")
		} else {
			return 200, nil
		}
	}
}

func UpdateForecast(forecastID int64, weather payload.UpdateDataForecast, db *sql.DB) (errCode int16, err error) {
	locationID, locationErrorCode, locationError := GetLocationID(weather.Country, weather.City, db)
	if locationID == nil {
		return locationErrorCode, locationError
	} else {
		row, err := db.Query("UPDATE forecast_weather_table SET location_id = $1, timestamp = $2, temperature_high = $3, temperature_low = $4, humidity = $5, weather_condition = $6 WHERE forecast_id = $7 RETURNING *;", locationID, weather.Timestamp, weather.TemperatureHigh, weather.TemperatureLow, weather.Humidity, weather.WeatherCondition, forecastID)
		if err != nil {
			return 500, errors.New("internal server error")
		}
		defer func(row *sql.Rows) {
			err := row.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(row)
		if !row.Next() {
			return 404, errors.New("weather not found")
		} else {
			return 200, nil
		}
	}
}
