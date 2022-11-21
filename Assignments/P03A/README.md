## P03A - Missile Command
### Dakota Wilson
### Description:

This project makes a table with the 70 largest military bases and also creates a table of random missile positions to launch towards them.

### Files

|   #   | File/folder              | Description                                                         |
| :---: | ------------------------ | ------------------------------------------------------------------- |
|   1   | LoadTables.go            | Main driver of my project that loads all the tables                 |
|   2   | SQL                      | SQL files that make the tables, adds indexes, and other useful sql  |
|   3   | Basesout.geojson         | Geojson of first 10 bases in the table                              |
|   4   | Missilesout.geojson      | Geojson of first 10 missiles in the table                           |


### Instructions

- Make sure you install library from `github.com/lib/pq`
- Make sure you install library from `github.com/paulmach/orb`
- Run `LoadTables.go` to create the output

- Example Command:
    - `go run LoadTables.go`