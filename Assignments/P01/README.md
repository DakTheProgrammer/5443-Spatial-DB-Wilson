## P01 - Project setup
### Dakota Wilson
### Description:

This project uses a postgres server to make an API about earthquakes. This is done through a Go program that uses a CSV data file to load its info into the database. After the data is all loaded there is a Go program that creates an API with 3 basic routes to ensure that the database is working. One of these routes uses a spatial query to assure the usage of PostGIS.

### Files

|   #   | File/folder              | Description                                        |
| :---: | ------------------------ | -------------------------------------------------- |
|   1   | EarthquakesAPI.go        | Main driver of my project that launches the API    |
|   2   | SQL                      | SQL files that make the table, test queries, etc.  |
|   3   | Load/LoadEarthquakes.go  | Helper go script that loads the database with data |
|   4   | Load/earthquakes.csv     | Data file with all of the earthquake data used     |
|   5   | IMG                      | Folder with images for presentation                |

### Instructions

- Make sure you install library from `github.com/lib/pq`
- Make sure you install library from `github.com/gin-gonic/gin`
- Make sure you install library from `github.com/paulmach/orb`
- Run `EarthquakesAPI.go` to run the API

- Example Command:
    - `go run EarthquakesAPI.go`
    - `go run Load/LoadEarthquakes.go`

# FindAll
<img src="IMG/All1.jpg">
<img src="IMG/All2.jpg">

# FindOne
### ID 13543
<img src="IMG/One1.jpg">

### ID 2598
<img src="IMG/One2.jpg">

# FindClosest
### Italy
<img src="IMG/Closest1.jpg">

### Ocean between Australia and South America
<img src="IMG/Closest2.jpg">

### Location out of range
<img src="IMG/Closest3.jpg">

### No Lon or Lat given
<img src="IMG/Closest4.jpg">