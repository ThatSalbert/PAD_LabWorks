package database

import (
	"database/sql"
	"disaster-service/payload"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func ConnectToDB(maxConnections int, DB_HOSTNAME string, DB_PORT string) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", "user=postgres password=postgres, dbname=disaster-service-database host="+DB_HOSTNAME+" port="+DB_PORT+" sslmode=disable")
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

func GetDisasterTypeID(disasterName string, db *sql.DB) (disasterTypeID *int, errCode int16, err error) {
	rows, err := db.Query("SELECT disaster_type_id FROM disaster_type_table WHERE UPPER(disaster_name) LIKE UPPER($1)", disasterName)
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
		return nil, 404, errors.New("disaster type not found")
	} else {
		if err = rows.Scan(&disasterTypeID); err != nil {
			return nil, 500, errors.New("internal server error")
		} else {
			return disasterTypeID, 200, nil
		}
	}
}

func GetDisasterTypes(db *sql.DB) (types []payload.DisasterType, errCode int16, err error) {
	rows, err := db.Query("SELECT disaster_name, disaster_description FROM disaster_type_table")
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
		var disasterType payload.DisasterType
		if err = rows.Scan(&disasterType.DisasterName, &disasterType.DisasterDescription); err != nil {
			return nil, 500, errors.New("internal server error")
		}
		if disasterType.DisasterName != "" {
			types = append(types, disasterType)
		} else {
			return nil, 404, errors.New("disaster types not found")
		}
	}
	return types, 200, nil
}

func GetDisasterList(db *sql.DB, country string, city string, active bool) (disasterList []payload.DisasterList, errCode int16, err error) {
	locationID, locationIDErrCode, locationIDErr := GetLocationID(country, city, db)
	if locationID == nil {
		return nil, locationIDErrCode, locationIDErr
	} else {
		if active {
			rows, err := db.Query("SELECT lt.country, lt.city, dtt.disaster_name, dlt.timestamp_start, dlt.timestamp_end, dlt.severity, dlt.description FROM disaster_list_table dlt INNER JOIN disaster_type_table dtt ON dtt.disaster_type_id = dlt.disaster_type_id INNER JOIN location_table lt ON lt.location_id = dlt.location_id WHERE dlt.location_id = $1 AND dlt.timestamp_end >= NOW()", locationID)
			if err != nil {
				return nil, 500, errors.New("internal server error")
			}
			defer func(rows *sql.Rows) {
				err := rows.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(rows)
			var disasterListElement payload.DisasterList
			for rows.Next() {
				var disaster payload.Disaster
				err := rows.Scan(&disasterListElement.Country, &disasterListElement.City, &disaster.DisasterName, &disaster.DisasterTimestampStart, &disaster.DisasterTimestampEnd, &disaster.DisasterSeverity, &disaster.DisasterDescription)
				if err != nil {
					return nil, 500, errors.New("internal server error")
				}
				disasterTypeID, disasterTypeIDErrCode, disasterTypeIDErr := GetDisasterTypeID(disaster.DisasterName, db)
				if disasterTypeID == nil {
					return nil, disasterTypeIDErrCode, disasterTypeIDErr
				} else {
					disasterListElement.Disasters = append(disasterListElement.Disasters, disaster)
				}
			}
			disasterList = append(disasterList, disasterListElement)
			if len(disasterListElement.Disasters) == 0 {
				return nil, 404, errors.New("no active disasters")
			} else {
				return disasterList, 200, nil
			}
		} else {
			rows, err := db.Query("SELECT lt.country, lt.city, dtt.disaster_name, dlt.timestamp_start, dlt.timestamp_end, dlt.severity, dlt.description FROM disaster_list_table dlt INNER JOIN disaster_type_table dtt ON dtt.disaster_type_id = dlt.disaster_type_id INNER JOIN location_table lt ON lt.location_id = dlt.location_id WHERE dlt.location_id = $1", locationID)
			if err != nil {
				return nil, 500, errors.New("internal server error")
			}
			defer func(rows *sql.Rows) {
				err := rows.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(rows)
			var disasterListElement payload.DisasterList
			for rows.Next() {
				var disaster payload.Disaster
				err := rows.Scan(&disasterListElement.Country, &disasterListElement.City, &disaster.DisasterName, &disaster.DisasterTimestampStart, &disaster.DisasterTimestampEnd, &disaster.DisasterSeverity, &disaster.DisasterDescription)
				if err != nil {
					return nil, 500, errors.New("internal server error")
				}
				disasterTypeID, disasterTypeIDErrCode, disasterTypeIDErr := GetDisasterTypeID(disaster.DisasterName, db)
				if disasterTypeID == nil {
					return nil, disasterTypeIDErrCode, disasterTypeIDErr
				} else {
					disasterListElement.Disasters = append(disasterListElement.Disasters, disaster)
				}
			}
			disasterList = append(disasterList, disasterListElement)
			if len(disasterListElement.Disasters) == 0 {
				return nil, 404, errors.New("no disasters")
			} else {
				return disasterList, 200, nil
			}
		}
	}
}

func AddAlert(alert payload.AddAlert, db *sql.DB) (alertID *int, errCode int16, err error) {
	locationID, locationIDErrCode, locationIDErr := GetLocationID(alert.Country, alert.City, db)
	disasterTypeID, disasterTypeIDErrCode, disasterTypeIDErr := GetDisasterTypeID(alert.DisasterName, db)
	if locationID == nil {
		return nil, locationIDErrCode, locationIDErr
	} else {
		if disasterTypeID == nil {
			return nil, disasterTypeIDErrCode, disasterTypeIDErr
		} else {
			rows, err := db.Query("INSERT INTO disaster_list_table (disaster_type_id, location_id, timestamp_start, timestamp_end, severity, description) SELECT $1, $2, $3, $4, CAST($5 AS VARCHAR), $6 WHERE NOT EXISTS (SELECT 1 FROM disaster_list_table WHERE disaster_type_id = $1 AND location_id = $2 AND timestamp_start = $3 AND timestamp_end = $4 AND severity = CAST($5 AS VARCHAR) AND description = $6) RETURNING disaster_id", disasterTypeID, locationID, alert.TimestampStart, alert.TimestampEnd, alert.Severity, alert.Description)
			if err != nil {
				fmt.Println(err)
				return nil, 500, errors.New("internal server error")
			}
			defer func(rows *sql.Rows) {
				err := rows.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(rows)
			if !rows.Next() {
				return nil, 409, errors.New("alert already exists")
			} else {
				if err = rows.Scan(&alertID); err != nil {
					return nil, 500, errors.New("internal server error")
				} else {
					return alertID, 200, nil
				}
			}
		}
	}
}

func UpdateAlert(alertID int, alert payload.UpdateAlert, db *sql.DB) (errCode int16, err error) {
	locationID, locationIDErrCode, locationIDErr := GetLocationID(alert.Country, alert.City, db)
	disasterTypeID, disasterTypeIDErrCode, disasterTypeIDErr := GetDisasterTypeID(alert.DisasterName, db)
	if locationID == nil {
		return locationIDErrCode, locationIDErr
	} else {
		if disasterTypeID == nil {
			return disasterTypeIDErrCode, disasterTypeIDErr
		} else {
			rows, err := db.Query("UPDATE disaster_list_table SET location_id = $1, disaster_type_id = $2, timestamp_start = $3, timestamp_end = $4, severity = $5, description = $6 WHERE disaster_id = $7 RETURNING *;", locationID, disasterTypeID, alert.TimestampStart, alert.TimestampEnd, alert.Severity, alert.Description, alertID)
			if err != nil {
				return 500, errors.New("internal server error")
			}
			defer func(rows *sql.Rows) {
				err := rows.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(rows)
			if !rows.Next() {
				return 404, errors.New("alert not found")
			} else {
				return 200, nil
			}
		}
	}
}
