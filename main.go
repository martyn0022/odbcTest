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

	rows, err := db.Query("select top 10 * from NYTaxi.Zones")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	var results []interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range columns {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			log.Fatal(err)
		}

		for _, v := range values {
			if b, ok := v.([]byte); ok {
				results = append(results, string(b))
			} else {
				results = append(results, v)
			}
		}
	}

	// Convert results to JSON
	jsonData, err := json.Marshal(results)
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
