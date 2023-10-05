package database

import (
	"database/sql"
	"errors"
	"fmt"
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
	return db, err
}

func GetCities(country string, db *sql.DB) (locations []payload.Location, err error) {
	rows, err := db.Query("SELECT city, longitude, latitude FROM location_table WHERE country=$1", country)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		locations = append(locations, location)
	}
	return locations, nil
}

func GetCurrentWeather(city string, db *sql.DB) (weather []payload.CurrentWeatherResponse, err error) {
	locationIDRow, err := db.Query("SELECT location_id FROM location_table WHERE UPPER(city) LIKE UPPER($1);", city)
	if err != nil {
		return nil, err
	}
	defer func(locationID *sql.Rows) {
		err := locationID.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(locationIDRow)
	var (
		locationID *int
	)
	for locationIDRow.Next() {
		if err = locationIDRow.Scan(&locationID); err != nil {
			return nil, err
		}
	}
	if locationID == nil {
		err := errors.New("location not found")
		return nil, err
	} else {
		rows, err := db.Query("SELECT DISTINCT ON (location_table.location_id) location_table.country, location_table.city, location_table.longitude, location_table.latitude, current_weather_table.timestamp, current_weather_table.temperature, current_weather_table.humidity, current_weather_table.weather_condition FROM current_weather_table INNER JOIN location_table ON location_table.location_id = $1 WHERE location_table.location_id = $1 AND current_weather_table.timestamp >= DATE_TRUNC('day', NOW()) AND current_weather_table.timestamp <= NOW() ORDER BY location_table.location_id, current_weather_table.timestamp DESC;", locationID)
		if err != nil {
			return nil, err
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(rows)
		for rows.Next() {
			var currentWeather payload.CurrentWeatherResponse
			if err = rows.Scan(&currentWeather.Country, &currentWeather.Location.City, &currentWeather.Location.Longitude, &currentWeather.Location.Latitude, &currentWeather.Timestamp, &currentWeather.Temperature, &currentWeather.Humidity, &currentWeather.WeatherCondition); err != nil {
				return nil, err
			}
			weather = append(weather, currentWeather)
		}
		return weather, nil
	}
}

func GetForecastWeather(city string, db *sql.DB) (forecast []payload.ForecastWeatherResponse, err error) {
	locationIDRow, err := db.Query("SELECT location_id FROM location_table WHERE UPPER(city) LIKE UPPER($1);", city)
	if err != nil {
		return nil, err
	}
	defer func(locationID *sql.Rows) {
		err := locationID.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(locationIDRow)
	var (
		locationID *int
	)
	for locationIDRow.Next() {
		if err = locationIDRow.Scan(&locationID); err != nil {
			return nil, err
		}
	}
	if locationID == nil {
		err := errors.New("location not found")
		return nil, err
	} else {
		rows, err := db.Query("SELECT lt.country, lt.city, lt.longitude, lt.latitude, fwt.timestamp, fwt.temperature_high, fwt.temperature_low, fwt.humidity, fwt.weather_condition FROM forecast_weather_table AS fwt INNER JOIN location_table AS lt ON lt.location_id = fwt.location_id WHERE fwt.location_id = $1 AND fwt.timestamp >= NOW() AND fwt.timestamp <= NOW()::DATE + INTERVAL '5 days' ORDER BY lt.location_id, fwt.timestamp;", locationID)
		if err != nil {
			log.Fatal(err)
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(rows)
		if errors.Is(rows.Scan(), sql.ErrNoRows) {
			fmt.Print("No rows found")
			return nil, nil
		}
		var forecastWeather payload.ForecastWeatherResponse
		for rows.Next() {
			var forecastWeatherDay payload.ForecastWeatherDay
			if err = rows.Scan(&forecastWeather.Country, &forecastWeather.Location.City, &forecastWeather.Location.Longitude, &forecastWeather.Location.Latitude, &forecastWeatherDay.Timestamp, &forecastWeatherDay.TemperatureHigh, &forecastWeatherDay.TemperatureLow, &forecastWeatherDay.Humidity, &forecastWeatherDay.WeatherCondition); err != nil {
				fmt.Print("Error called: + ", err)
				return nil, err
			}
			forecastWeather.ForecastWeatherDays = append(forecastWeather.ForecastWeatherDays, forecastWeatherDay)
		}
		if len(forecastWeather.ForecastWeatherDays) == 0 {
			return nil, errors.New("no forecast weather found")
		} else {
			forecast = append(forecast, forecastWeather)
			return forecast, nil
		}
	}
}

func AddCurrentWeather(weather payload.AddDataWeather, db *sql.DB) (err error) {
	locationIDRow, err := db.Query("SELECT location_id FROM location_table WHERE UPPER(city) LIKE UPPER($1);", weather.City)
	if err != nil {
		return err
	}
	defer func(locationID *sql.Rows) {
		err := locationID.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(locationIDRow)
	var (
		locationID *int
	)
	for locationIDRow.Next() {
		if err = locationIDRow.Scan(&locationID); err != nil {
			return err
		} else {
			valueExistsRows, err := db.Query("SELECT EXISTS(SELECT 1 FROM current_weather_table WHERE location_id = $1 AND timestamp = $2);", locationID, weather.Timestamp)
			if err != nil {
				return err
			}
			defer func(valueExists *sql.Rows) {
				err := valueExists.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(valueExistsRows)
			var (
				valueExists *bool
			)
			for valueExistsRows.Next() {
				if err = valueExistsRows.Scan(&valueExists); err != nil {
					return err
				} else {
					if *valueExists == false {
						err = db.QueryRow("INSERT INTO current_weather_table (location_id, timestamp, temperature, humidity, weather_condition) VALUES ((SELECT location_id FROM location_table WHERE UPPER(city) LIKE UPPER($1)), $2, $3, $4, $5);", weather.City, weather.Timestamp, weather.Temperature, weather.Humidity, weather.WeatherCondition).Scan()
						if errors.Is(err, sql.ErrNoRows) {
							return nil
						} else {
							return err
						}
					} else {
						return errors.New("value already exists")
					}
				}
			}
		}
	}
	return nil
}

func AddForecastWeather(weather payload.AddDataForecast, db *sql.DB) (err error) {
	locationIDRow, err := db.Query("SELECT location_id FROM location_table WHERE UPPER(city) LIKE UPPER($1);", weather.City)
	if err != nil {
		return err
	}
	defer func(locationID *sql.Rows) {
		err := locationID.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(locationIDRow)
	var (
		locationID *int
	)
	for locationIDRow.Next() {
		if err = locationIDRow.Scan(&locationID); err != nil {
			return err
		} else {
			valueExistsRows, err := db.Query("SELECT EXISTS(SELECT 1 FROM forecast_weather_table WHERE location_id = $1 AND timestamp = $2);", locationID, weather.Timestamp)
			if err != nil {
				return err
			}
			defer func(valueExists *sql.Rows) {
				err := valueExists.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(valueExistsRows)
			var (
				valueExists *bool
			)
			for valueExistsRows.Next() {
				if err = valueExistsRows.Scan(&valueExists); err != nil {
					return err
				} else {
					if *valueExists == false {
						err = db.QueryRow("INSERT INTO forecast_weather_table (location_id, timestamp, temperature_high, temperature_low, humidity, weather_condition) VALUES ((SELECT location_id FROM location_table WHERE UPPER(city) LIKE UPPER($1)), $2, $3, $4, $5, $6);", weather.City, weather.Timestamp, weather.TemperatureHigh, weather.TemperatureLow, weather.Humidity, weather.WeatherCondition).Scan()
						if errors.Is(err, sql.ErrNoRows) {
							return nil
						} else {
							return err
						}
					} else {
						return errors.New("value already exists")
					}
				}
			}
		}
	}
	return nil
}

func DeleteWeather(city string, timestamp string, db *sql.DB) (err error) {
	locationIDRow, err := db.Query("SELECT location_id FROM location_table WHERE UPPER(city) LIKE UPPER($1);", city)
	if err != nil {
		return err
	}
	defer func(locationID *sql.Rows) {
		err := locationID.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(locationIDRow)
	var (
		locationID *int
	)
	for locationIDRow.Next() {
		if err = locationIDRow.Scan(&locationID); err != nil {
			return err
		} else {
			valueExistsRows, err := db.Query("SELECT EXISTS(SELECT 1 FROM current_weather_table WHERE location_id = $1 AND timestamp = $2);", locationID, timestamp)
			if err != nil {
				return err
			}
			defer func(valueExists *sql.Rows) {
				err := valueExists.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(valueExistsRows)
			var (
				valueExists *bool
			)
			for valueExistsRows.Next() {
				if err = valueExistsRows.Scan(&valueExists); err != nil {
					return err
				} else {
					if *valueExists == true {
						err = db.QueryRow("DELETE FROM current_weather_table WHERE location_id = $1 AND timestamp = $2;", locationID, timestamp).Scan()
						if errors.Is(err, sql.ErrNoRows) {
							return nil
						} else {
							return err
						}
					} else {
						return errors.New("value does not exist")
					}
				}
			}
		}
	}
	return nil
}

func DeleteForecast(city string, timestamp string, db *sql.DB) (err error) {
	locationIDRow, err := db.Query("SELECT location_id FROM location_table WHERE UPPER(city) LIKE UPPER($1);", city)
	if err != nil {
		return err
	}
	defer func(locationID *sql.Rows) {
		err := locationID.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(locationIDRow)
	var (
		locationID *int
	)
	for locationIDRow.Next() {
		if err = locationIDRow.Scan(&locationID); err != nil {
			return err
		} else {
			valueExistsRows, err := db.Query("SELECT EXISTS(SELECT 1 FROM forecast_weather_table WHERE location_id = $1 AND timestamp = $2);", locationID, timestamp)
			if err != nil {
				return err
			}
			defer func(valueExists *sql.Rows) {
				err := valueExists.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(valueExistsRows)
			var (
				valueExists *bool
			)
			for valueExistsRows.Next() {
				if err = valueExistsRows.Scan(&valueExists); err != nil {
					return err
				} else {
					if *valueExists == true {
						err = db.QueryRow("DELETE FROM forecast_weather_table WHERE location_id = $1 AND timestamp = $2;", locationID, timestamp).Scan()
						if errors.Is(err, sql.ErrNoRows) {
							return nil
						} else {
							return err
						}
					} else {
						return errors.New("value does not exist")
					}
				}
			}
		}
	}
	return nil
}
