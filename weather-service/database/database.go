package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func ConnectToDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "user=postgres password=1234567890 dbname=weather-service-database sslmode=disable")
	if db == nil {
		return nil, err
	}
	return db, err
}

func GetTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var tables []string
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}
