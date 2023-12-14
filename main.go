package main

import (
	"database/sql"
	"fmt"
	"log"

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

	fmt.Println("Fetched Results:")
	for i := 0; i < len(results); i += len(columns) {
		for j := 0; j < len(columns); j++ {
			fmt.Printf("%s: %v\t", columns[j], results[i+j])
		}
		fmt.Println()
	}
}
