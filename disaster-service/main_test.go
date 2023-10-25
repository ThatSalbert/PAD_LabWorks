package main

import (
	"database/sql"
	"disaster-service/database"
	"disaster-service/payload"
	_ "github.com/lib/pq"
	"log"
	"testing"
)

func dbConnect() (db *sql.DB) {
	db, _ = sql.Open("postgres", "user=postgres password=postgres dbname=test-service-database host=test-database port=5432 sslmode=disable")
	return db
}

func TestGetDisasterTypesList(t *testing.T) {
	db := dbConnect()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	_, statusCode, _ := database.GetDisasterTypes(db)

	if statusCode != 200 {
		t.Errorf("Expected status code 200, got %d", statusCode)
	}
}

func TestGetLocationID(t *testing.T) {
	db := dbConnect()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	locationID, statusCode, _ := database.GetLocationID("Moldova", "Balti", db)

	if statusCode != 200 && *locationID != 3 {
		t.Errorf("Expected status code 200 and location ID 3, got %d and location ID %d", statusCode, locationID)
	}

	locationID, statusCode, _ = database.GetLocationID("Moldova", "Chisinau", db)

	if statusCode != 200 && *locationID != 1 {
		t.Errorf("Expected status code 200 and location ID 1, got %d and location ID %d", statusCode, locationID)
	}

	locationID, statusCode, _ = database.GetLocationID("Moldova", "Cahul", db)

	if statusCode != 200 && *locationID != 6 {
		t.Errorf("Expected status code 200 and location ID 6, got %d and location ID %d", statusCode, locationID)
	}
}

func TestAddAlert(t *testing.T) {
	db := dbConnect()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	addAlertRequest := payload.AddAlert{
		Country:        "Moldova",
		City:           "Balti",
		DisasterName:   "Flood",
		TimestampStart: "2023-10-11 00:00:00",
		TimestampEnd:   "2023-11-03 00:00:00",
		Severity:       "High",
		Description:    "Flood in Balti",
	}

	_, statusCode, _ := database.AddAlert(addAlertRequest, db)

	if statusCode != 200 {
		t.Errorf("Expected status code 200, got %d", statusCode)
	}

	addAlertRequest = payload.AddAlert{
		Country:        "Moldova",
		City:           "Chisinau",
		DisasterName:   "nonexistentdisaster",
		TimestampStart: "2023-10-11 00:00:00",
		TimestampEnd:   "2023-11-03 00:00:00",
		Severity:       "High",
		Description:    "What",
	}

	_, statusCode, _ = database.AddAlert(addAlertRequest, db)

	if statusCode != 404 {
		t.Errorf("Expected status code 404, got %d", statusCode)
	}
}

func TestGetAlert(t *testing.T) {
	db := dbConnect()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	alert, statusCode, _ := database.GetDisasterList(db, "Moldova", "Balti", true)

	if statusCode != 200 && len(alert) != 1 {
		t.Errorf("Expected status code 200 and non-empty alert list, got %d and empty alert list", statusCode)
	}
}

func TestGetDisasterList(t *testing.T) {
	db := dbConnect()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	addAlertRequest := payload.AddAlert{
		Country:        "Moldova",
		City:           "Chisinau",
		DisasterName:   "Earthquake",
		TimestampStart: "2023-10-12 00:00:00",
		TimestampEnd:   "2023-11-05 00:00:00",
		Severity:       "High",
		Description:    "Earthquake in Chisinau x2",
	}

	database.AddAlert(addAlertRequest, db)

	_, statusCode, _ := database.GetDisasterList(db, "Moldova", "Chisinau", false)

	if statusCode != 200 {
		t.Errorf("Expected status code 200, got %d", statusCode)
	}
}
