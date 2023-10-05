package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

func ConnectToDB(maxConnections int) (*sql.DB, error) {
	db, err := sql.Open("postgres", "user=postgres password=1234567890 dbname=disaster-service-database sslmode=disable")
	if db == nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxConnections)
	db.SetMaxIdleConns(maxConnections / 2)
	db.SetConnMaxLifetime(time.Second * 5)
	return db, nil
}
