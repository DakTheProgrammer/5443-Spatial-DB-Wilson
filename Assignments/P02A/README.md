## P02A - World Data
### Dakota Wilson
### Description:

This project loads shape files and a csv into a postgres server and creates indexes on those tables. After this I output the first 10 rows of each table to a JSON file.

### Files

|   #   | File/folder              | Description                                                         |
| :---: | ------------------------ | ------------------------------------------------------------------- |
|   1   | MakeOutput.go            | Main driver of my project that outputs first 10 rows to a JSON file |
|   2   | SQL                      | SQL files that make the tables and adds indexes                     |
|   3   | Data/Military            | Has the shape file for military bases                               |
|   4   | Data/Rails               | Has the shape file for railroads                                    |
|   5   | Data/Roads               | Has the shape file for primary roads                                |
|   6   | Data/States              | Has the shape file for states                                       |
|   7   | Data/Timezones           | Has the shape file for timezones                                    |
|   8   | Data/airports.csv        | The csv file for airport data                                       |
|   9   | Data/LoadAirports.go     | Go script that loads all the airport data to Postgres               |


### Instructions

- Make sure you install library from `github.com/lib/pq`
- Make sure you install library from `github.com/paulmach/orb`
- Run `MakeOutput.go` to create the output

- Example Command:
    - `go run MakeOutput.go`
    - `go run Data/LoadAirports.go`