package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/ewkb"
)

// struct for columns of earthquake table
// NOTE: TO PASS TO API KEY MUST START WITH CAPITAL #ILOVEGO
type Earthquake struct {
	Id          int
	DateTime    string
	Location    orb.Point
	Depth       float32
	Magnitude   float32
	Calculation string
	Network     string
	Place       string
	Cause       string
}

// global variables that are reused multiple times
var (
	db          *sql.DB
	id          int
	dateTime    string
	location    orb.Point
	depth       float32
	magnitude   float32
	calculation string
	network     string
	place       string
	cause       string
)

// function that checks all errors (for reusability)
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// called before main
func init() {
	var err error

	jsonFile, err := os.Open("login.json")

	checkError(err)

	//reads entire json file and gets the byte total
	byteValue, err := ioutil.ReadAll(jsonFile)

	checkError(err)

	//gets all the file info and stores it into a map for json like structure
	var login map[string]string
	json.Unmarshal(byteValue, &login)
	jsonFile.Close()

	//creates connection to postgres
	db, err = sql.Open("postgres", "host="+login["host"]+" port="+login["port"]+" user="+login["user"]+
		" password="+login["password"]+" dbname="+login["dbname"]+" sslmode=disable")

	checkError(err)

	//pings DB to assure connection is working
	err = db.Ping()
	checkError(err)

	fmt.Println("Successfully connected to Postgres server!")
}

func main() {
	router := gin.Default()

	router.GET("/findAll", findAll)
	router.GET("/findOne/:id", findOne)
	router.GET("/findClosest", findClosest)

	router.Run("localhost:8080")

	db.Close()
}

// GET call that returns entire table
func findAll(c *gin.Context) {
	var earthquakes []Earthquake

	rows, err := db.Query("SELECT * FROM public.earthquakes;")

	checkError(err)

	//loops through all rows returned and stores them in each given variable to load into array
	for rows.Next() {
		rows.Scan(&id, &dateTime, ewkb.Scanner(&location), &depth, &magnitude, &calculation, &network,
			&place, &cause)

		//appends the array with each row of earthquake data
		earthquakes = append(earthquakes, Earthquake{Id: id, DateTime: dateTime, Location: location,
			Depth: depth, Magnitude: magnitude, Calculation: calculation, Network: network,
			Place: place, Cause: cause})
	}

	c.IndentedJSON(http.StatusOK, earthquakes)
}

// GET route that finds earthquake based on id#
func findOne(c *gin.Context) {
	ID := c.Param("id")

	rows, err := db.Query("SELECT * FROM public.earthquakes WHERE id = " + ID + ";")

	checkError(err)

	//see findClosest for documentation
	rows.Next()

	rows.Scan(&id, &dateTime, ewkb.Scanner(&location), &depth, &magnitude, &calculation, &network,
		&place, &cause)

	c.IndentedJSON(http.StatusOK, Earthquake{Id: id, DateTime: dateTime, Location: location,
		Depth: depth, Magnitude: magnitude, Calculation: calculation, Network: network,
		Place: place, Cause: cause})
}

// GET route that finds the closest earthquake to a given lon lat
func findClosest(c *gin.Context) {
	//uses query string to get params
	lon := c.Query("lon")
	lat := c.Query("lat")

	//converts strings to floats and assures its within SRID 4326(180 -> -180)
	longitude, err := strconv.ParseFloat(lon, 64)

	if err != nil || longitude > 180 || longitude < -180 {
		c.IndentedJSON(http.StatusOK, "Invalid longitude or latitude")
		return
	}

	latitude, err := strconv.ParseFloat(lat, 64)

	if err != nil || latitude > 180 || latitude < -180 {
		c.IndentedJSON(http.StatusOK, "Invalid longitude or latitude")
		return
	}

	rows, err := db.Query("SELECT *, ST_Distance('SRID=4326;POINT(" +
		lon + " " + lat + ")'::geometry, location) AS dist FROM public.earthquakes ORDER BY dist LIMIT 1;")

	checkError(err)

	var distance float64

	//Goes to row returned
	rows.Next()

	//Assigns row data to variables
	rows.Scan(&id, &dateTime, ewkb.Scanner(&location), &depth, &magnitude, &calculation, &network,
		&place, &cause, &distance)

	//Structure of JSON return
	type WithDistance struct {
		Earthquake Earthquake
		Distance   float64
	}

	c.IndentedJSON(http.StatusOK, WithDistance{Earthquake: Earthquake{Id: id, DateTime: dateTime, Location: location,
		Depth: depth, Magnitude: magnitude, Calculation: calculation, Network: network,
		Place: place, Cause: cause}, Distance: distance})
}
