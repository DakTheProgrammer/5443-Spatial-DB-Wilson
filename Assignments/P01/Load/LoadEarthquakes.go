package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	_ "github.com/lib/pq"
)

// function that checks all errors (for reusability)
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// global database connection
var db *sql.DB

// run before main
// see EarthquakeAPI.go for this functions documentation
func init() {
	var err error

	jsonFile, err := os.Open("../login.json")

	checkError(err)

	byteValue, err := ioutil.ReadAll(jsonFile)

	checkError(err)

	var result map[string]string
	json.Unmarshal(byteValue, &result)
	jsonFile.Close()

	db, err = sql.Open("postgres", "host="+result["host"]+" port="+result["port"]+" user="+result["user"]+
		" password="+result["password"]+" dbname="+result["dbname"]+" sslmode=disable")

	checkError(err)

	checkError(err)

	fmt.Println("Successfully connected to Postgres server!")
}

func main() {
	earthquakesCSV, err := os.Open("Earthquakes.csv")

	checkError(err)

	reader := csv.NewReader(earthquakesCSV)

	var sql string = "INSERT INTO public.earthquakes(id, datetime, location, depth, magnitude, calculation, network, place, cause) " +
		"VALUES (%s, '%s', ST_GeomFromText('POINT(%s %s)', 4326), %s, %s, '%s', '%s', '%s', '%s');"

	//used to read the leading row of CSV since unused
	reader.Read()

	for {
		earthquake, err := reader.Read()

		//reads until end of file
		if err == io.EOF {
			break
		}

		checkError(err)

		//loads each row into the database following given formatted string
		db.Exec(fmt.Sprintf(sql, earthquake[0], earthquake[1], earthquake[3], earthquake[2], earthquake[4], earthquake[5], earthquake[6], earthquake[7], earthquake[8], earthquake[9]))
		fmt.Print(earthquake[0], "\n")
	}

	db.Close()
	earthquakesCSV.Close()
}
