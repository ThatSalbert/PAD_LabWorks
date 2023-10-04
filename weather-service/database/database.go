package database

import (
	"database/sql"
	_ "github.com/lib/pq"
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
	db.SetConnMaxLifetime(time.Second * 10)
	return db, err
}

func GetCities(country string, db *sql.DB) (locations []payload.Location, err error) {
	rows, err := db.Query("SELECT city, longitude, latitude FROM location_table WHERE country=$1", country)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var location payload.Location
		if err = rows.Scan(&location.City, &location.Longitude, &location.Latitude); err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}
	return locations, nil
}
