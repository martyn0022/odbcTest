package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/alexbrainman/odbc"
)

func main() {
	db, err := sql.Open("odbc", "DSN=IRIS")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select TOP 10 VendorID, $PIECE(CAST(tpep_pickup_datetime as VARCHAR),' ',1), passenger_count, trip_distance from NYTaxi.Rides")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	var results [][]interface{} // Modified to hold arrays of values

	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range columns {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			log.Fatal(err)
		}

		results = append(results, values) // Append each row's values as an array
	}

	// Optional: Include column names in the JSON
	includeColumnNames := true // Change to false to exclude column names

	var finalResult interface{}
	if includeColumnNames {
		finalResult = map[string]interface{}{
			"columns": columns,
			"rows":    results,
		}
	} else {
		finalResult = results
	}

	// Convert finalResult to JSON
	jsonData, err := json.Marshal(finalResult)
	if err != nil {
		log.Fatal(err)
	}

	// Write JSON data to a file
	file, err := os.Create("results.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Results written to results.json")
}
