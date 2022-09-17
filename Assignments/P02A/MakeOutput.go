package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	_ "github.com/lib/pq"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/ewkb"
)

// global variable that are reused multiple times
var db *sql.DB

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

	//gets all the file info and stores it sql.NullInt64o a map for json like structure
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
	//reused variables throughout the program
	var (
		rows    *sql.Rows
		err     error
		columns []string

		i1  sql.NullInt64
		i2  sql.NullInt64
		i3  sql.NullInt64
		i4  sql.NullInt64
		i5  sql.NullInt64
		i6  sql.NullInt64
		s1  sql.NullString
		s2  sql.NullString
		s3  sql.NullString
		s4  sql.NullString
		s5  sql.NullString
		s6  sql.NullString
		s7  sql.NullString
		s8  sql.NullString
		s9  sql.NullString
		s10 sql.NullString
		s11 sql.NullString
		s12 sql.NullString
		d1  sql.NullFloat64
		d2  sql.NullFloat64
		p1  orb.Point
		mp1 orb.MultiPolygon
		ml1 orb.MultiLineString
	)

	//creates file if not created, clears the file, then prepares it for writing
	outfile, err := os.OpenFile("out.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	outfile.WriteString("{\"Airports\":[\n")

	//sql used in this program
	var AirportSQL string = "SELECT * FROM public.airports WHERE id <= 10"
	var MilBasesSQL string = "SELECT * FROM public.militarybases WHERE gid <= 10"
	var RoadsSQL string = "SELECT * FROM public.primaryroads WHERE gid <= 10"
	var RailroadsSQL string = "SELECT * FROM public.railroads WHERE gid <= 10"
	var StatesSQL string = "SELECT * FROM public.states WHERE gid <= 10"
	var TimezoneSQL string = "SELECT * FROM public.timezones WHERE gid <= 10"

	//map for JSON like output
	var m = make(map[string]any)

	rows, err = db.Query(AirportSQL)

	checkError(err)

	//gets the columns of the table as keys for output
	columns, err = rows.Columns()

	checkError(err)

	for rows.Next() {
		//loads variables for output
		rows.Scan(&i1, &s1, &s2, &s3, &s4, &s5, ewkb.Scanner(&p1), &i2, &i3, &s6, &s7, &s8)

		//stores all variables with correct keys
		m[columns[0]] = i1.Int64
		m[columns[1]] = s1.String
		m[columns[2]] = s2.String
		m[columns[3]] = s3.String
		m[columns[4]] = s4.String
		m[columns[5]] = s5.String
		m[columns[6]] = p1
		m[columns[7]] = i2.Int64
		m[columns[8]] = i3.Int64
		m[columns[9]] = s6.String
		m[columns[10]] = s7.String
		m[columns[11]] = s8.String

		//makes the output JSON like
		out, err := json.Marshal(m)

		checkError(err)

		//writes the row to the file
		outfile.WriteString(string(out) + ",\n")
	}

	//removes extra ","
	outfile.Seek(-2, io.SeekEnd)
	outfile.WriteString("\n")

	//----------------------------------------------------
	//----------------------------------------------------
	//----------------------------------------------------
	// ^see above code comments since this is all repeat^
	//----------------------------------------------------
	//----------------------------------------------------
	//----------------------------------------------------
	outfile.WriteString("],\n\"MilitaryBases\": [\n")

	m = make(map[string]any)

	rows, err = db.Query(MilBasesSQL)

	checkError(err)

	columns, err = rows.Columns()

	checkError(err)

	for rows.Next() {
		rows.Scan(&i1, &s1, &s2, &s3, &s4, &d1, &d2, &s5, &s6, ewkb.Scanner(&mp1))

		m[columns[0]] = i1.Int64
		m[columns[1]] = s1.String
		m[columns[2]] = s2.String
		m[columns[3]] = s3.String
		m[columns[4]] = s4.String
		m[columns[5]] = d1.Float64
		m[columns[6]] = d2.Float64
		m[columns[7]] = s5.String
		m[columns[8]] = s6.String
		m[columns[9]] = mp1

		out, err := json.Marshal(m)

		checkError(err)

		outfile.WriteString(string(out) + ",\n")
	}

	outfile.Seek(-2, io.SeekEnd)
	outfile.WriteString("\n")

	outfile.WriteString("],\n\"PrimaryRoads\": [\n")

	m = make(map[string]any)

	rows, err = db.Query(RoadsSQL)

	checkError(err)

	columns, err = rows.Columns()

	checkError(err)

	for rows.Next() {
		rows.Scan(&i1, &s1, &s2, &s3, &s4, ewkb.Scanner(&ml1))

		m[columns[0]] = i1.Int64
		m[columns[1]] = s1.String
		m[columns[2]] = s2.String
		m[columns[3]] = s3.String
		m[columns[4]] = s4.String
		m[columns[5]] = ml1

		out, err := json.Marshal(m)

		checkError(err)

		outfile.WriteString(string(out) + ",\n")
	}

	outfile.Seek(-2, io.SeekEnd)
	outfile.WriteString("\n")

	outfile.WriteString("],\n\"RailRoads\": [\n")

	m = make(map[string]any)

	rows, err = db.Query(RailroadsSQL)

	checkError(err)

	columns, err = rows.Columns()

	checkError(err)

	for rows.Next() {
		rows.Scan(&i1, &s1, &s2, &s3, ewkb.Scanner(&ml1))

		m[columns[0]] = i1.Int64
		m[columns[1]] = s1.String
		m[columns[2]] = s2.String
		m[columns[3]] = s3.String
		m[columns[4]] = ml1

		out, err := json.Marshal(m)

		checkError(err)

		outfile.WriteString(string(out) + ",\n")
	}

	outfile.Seek(-2, io.SeekEnd)
	outfile.WriteString("\n")

	outfile.WriteString("],\n\"States\": [\n")

	m = make(map[string]any)

	rows, err = db.Query(StatesSQL)

	checkError(err)

	columns, err = rows.Columns()

	checkError(err)

	for rows.Next() {
		rows.Scan(&i1, &s1, &s2, &s3, &s4, &s5, &s6, &s7, &s8, &s9, &s10, &d1, &d2, &s11, &s12, &mp1)

		m[columns[0]] = i1.Int64
		m[columns[1]] = s1.String
		m[columns[2]] = s2.String
		m[columns[3]] = s3.String
		m[columns[4]] = s4.String
		m[columns[5]] = s5.String
		m[columns[6]] = s6.String
		m[columns[7]] = s7.String
		m[columns[8]] = s8.String
		m[columns[9]] = s9.String
		m[columns[10]] = s10.String
		m[columns[11]] = d1.Float64
		m[columns[12]] = d2.Float64
		m[columns[13]] = s11.String
		m[columns[14]] = s12.String
		m[columns[15]] = mp1

		out, err := json.Marshal(m)

		checkError(err)

		outfile.WriteString(string(out) + ",\n")
	}

	outfile.Seek(-2, io.SeekEnd)
	outfile.WriteString("\n")

	outfile.WriteString("],\n\"TimeZones\": [\n")

	m = make(map[string]any)

	rows, err = db.Query(TimezoneSQL)

	checkError(err)

	columns, err = rows.Columns()

	checkError(err)

	for rows.Next() {
		rows.Scan(&i1, &i2, &i3, &s1, &s2, &i4, &i5, &s3, &d1, &s4, &s5, &s6, &s7, &s8, &s9, &i6, &mp1)

		m[columns[0]] = i1.Int64
		m[columns[1]] = i2.Int64
		m[columns[2]] = i3.Int64
		m[columns[3]] = s1.String
		m[columns[4]] = s2.String
		m[columns[5]] = i4.Int64
		m[columns[6]] = i5.Int64
		m[columns[7]] = s3.String
		m[columns[8]] = d1.Float64
		m[columns[9]] = s4.String
		m[columns[10]] = s5.String
		m[columns[11]] = s6.String
		m[columns[12]] = s7.String
		m[columns[13]] = s8.String
		m[columns[14]] = s9.String
		m[columns[15]] = i6.Int64
		m[columns[16]] = mp1

		out, err := json.Marshal(m)

		checkError(err)

		outfile.WriteString(string(out) + ",\n")
	}

	outfile.Seek(-2, io.SeekEnd)
	outfile.WriteString("\n")

	outfile.WriteString("]}")

	outfile.Close()
	db.Close()
}
