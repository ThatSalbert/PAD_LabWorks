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

func GetCurrentWeather(city string, db *sql.DB) (weather []payload.CurrentWeather, err error) {
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
		rows, err := db.Query("SELECT location_table.country, location_table.city, location_table.longitude, location_table.latitude, current_weather_table.timestamp, current_weather_table.temperature, current_weather_table.humidity, current_weather_table.weather_condition FROM current_weather_table INNER JOIN location_table ON location_table.location_id = $1;", locationID)
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
			var currentWeather payload.CurrentWeather
			if err = rows.Scan(&currentWeather.Country, &currentWeather.Location.City, &currentWeather.Location.Longitude, &currentWeather.Location.Latitude, &currentWeather.Timestamp, &currentWeather.Temperature, &currentWeather.Humidity, &currentWeather.WeatherCondition); err != nil {
				return nil, err
			}
			weather = append(weather, currentWeather)
		}
		return weather, nil
	}
}
