package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"

	_ "github.com/lib/pq"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/ewkb"
	"github.com/paulmach/orb/geojson"
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

	jsonFile, err := os.Open("../login.json")

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

	var BasesSQL string = "SELECT * from (SELECT ST_X(Center::geometry) as x, ST_Y(Center::geometry) as y, gid, (landarea + waterarea)::BIGINT as area from (SELECT ST_Centroid(miltable.geog) as Center, miltable.gid as gid, miltable.aland as landarea, miltable.awater as waterarea from militarybases as miltable)as useless) as why where x < -66 and x > -125 and y > 23 and y < 50 ORDER BY area DESC limit 70;"
	var InsertBasesSQL string = "INSERT INTO public.gamebases(lid, gid, center, health, missilecount, area) VALUES ($1, $2, ST_GeomFromText('POINT(%f %f)', 4326), $3, $4, $5);"
	var MissileStartsSQL string = "SELECT ST_AsGeoJSON(J.*) FROM (SELECT ST_GeneratePoints(geom, 800, 1996) FROM (SELECT ST_Buffer(ST_GeomFromText('LINESTRING(-129.7844079 19.7433195,-61.9513812 19.7433195 , -61.9513812 54.3457868,-129.7844079 54.3457868, -129.7844079 19.7433195)'), 1, 'endcap=round join=round')  As geom ) as s )as J;"
	var InsertMissilesSQL string = "INSERT INTO public.missiles(mid, start, speed, angle, altitude, blastradius) VALUES ($1, ST_GeomFromText('POINT(%f %f)', 4326), $2, $3, $4, $5);"
	var MissileAngleSQL string = "SELECT degrees(ST_Azimuth( ST_Point(%f, %f),  ST_Point(%f, %f)))"
	var RandomBaseSQL string = "SELECT center FROM public.gamebases WHERE lid = $1"

	var (
		err  error
		rows *sql.Rows

		x    float64
		y    float64
		gid  int
		area int64

		json   string
		points orb.MultiPoint
		point  orb.Point
		ang    float64
	)

	rows, err = db.Query(BasesSQL)

	checkError(err)

	var l int = 1

	for rows.Next() {
		rows.Scan(&x, &y, &gid, &area)

		//uses base area to set a health value and missile storage
		fmt.Println(db.Exec(fmt.Sprintf(InsertBasesSQL, x, y), l, gid, (area / 100000000 * 100), (area / 100000000 * 80), area))

		l++
	}

	rows, err = db.Query(MissileStartsSQL)

	checkError(err)

	for rows.Next() {
		rows.Scan(&json)

		//had to convert to geojson to allocate memory for multipoint
		rawjson := []byte(json)

		featureCollection, err := geojson.UnmarshalFeature(rawjson)

		checkError(err)

		//goes from geometry type to multipoint
		points = featureCollection.Geometry.(orb.MultiPoint)

		for i := 0; i < len(points); i++ {
			//targets a random base
			base, err := db.Query(RandomBaseSQL, rand.Intn(69)+1)

			checkError(err)

			for base.Next() {
				base.Scan(ewkb.Scanner(&point))

				//uses points to get angle missile should go from origin
				miss, err := db.Query(fmt.Sprintf(MissileAngleSQL, points[i].X(), points[i].Y(), point.X(), point.Y()))

				checkError(err)

				for miss.Next() {
					miss.Scan(&ang)

					//uses random values for speed, altitude and blast radius
					fmt.Println(db.Exec(fmt.Sprintf(InsertMissilesSQL, points[i].X(), points[i].Y()), i, 20000+rand.Float64()*(25000-20000), ang, rand.Intn(15000)+15000, 100+rand.Intn(400)))
				}

			}

		}
	}
}
