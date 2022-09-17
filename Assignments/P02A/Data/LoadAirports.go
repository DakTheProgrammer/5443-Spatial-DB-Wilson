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

// function that checks if a string should be NULL type for SQL
func CheckNullString(s string) sql.NullString {
	if len(s) != 0 {
		return sql.NullString{
			String: s,
			Valid:  true,
		}
	} else {
		return sql.NullString{}
	}
}

// global database connection
var db *sql.DB

// run before main
// see earthquakesAPI.go for this functions documentation
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
	airportsCSV, err := os.Open("airports.csv")

	checkError(err)

	reader := csv.NewReader(airportsCSV)

	var ins string = "INSERT INTO public.airports(id, name, city, country, shortcode, longcode, location, elevation, gmt, tz_short, timezone, type) " +
		"VALUES (%s, '%s', '%s', '%s', $1, '%s', ST_GeomFromText('POINT(%s %s)', 4326), %s, $2, $3, $4, '%s');"

	//used to read the leading row of CSV since unused
	reader.Read()

	for {
		airports, err := reader.Read()

		//reads until end of file
		if err == io.EOF {
			break
		}

		checkError(err)

		for i := 0; i < 12; i++ {
			if airports[i] == "\\N" {
				airports[i] = ""
			}
		}

		//loads each row into the database following given formatted string
		_, err = db.Exec(fmt.Sprintf(ins, airports[0], airports[1], airports[2], airports[3], airports[5], airports[6], airports[7], airports[8], airports[12]),
			CheckNullString(airports[4]), CheckNullString(airports[9]), CheckNullString(airports[10]), CheckNullString(airports[11]))

		checkError(err)

		fmt.Print(airports[0], "\n")

	}

	db.Close()
	airportsCSV.Close()
}
