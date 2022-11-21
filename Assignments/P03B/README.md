## P03B - Missile Command
### Dakota Wilson
### Description:

This project takes in API calls to make PSQL tables around missile command. With that it also uses the information there is then calculations(PostGIS) made to determine missile intercepts over a region to shoot down. With these calculations there is then a POST request made to the API to ensure missiles where intercepted.

### Files

|   #   | File/folder              | Description                                                         |
| :---: | ------------------------ | ------------------------------------------------------------------- |
|   1   | Solution.py              | Main driver of my project that does all the API requests and tables |
|   2   | SQL                      | SQL files that make the tables, adds indexes, and other useful sql  |
|   3   | SampleAPICalls           | JSON files used for testing missiles                                |


### Instructions

- Make sure you install library from `pip install geojson`
- Make sure you install library from `pip install psycopg2`

- Example Command:
    - `python Solution.py`