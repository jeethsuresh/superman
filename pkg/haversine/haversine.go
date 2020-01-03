package haversine

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeethsuresh/superman/pkg/geolocation"
	"github.com/jeethsuresh/superman/pkg/sanitize"

	//Well-known SQLite3 driver
	_ "github.com/mattn/go-sqlite3"
)

const currentGeo = "currentGeo"
const travelToCurrentGeoSuspicious = "travelToCurrentGeoSuspicious"
const travelFromCurrentGeoSuspicious = "travelFromCurrentGeoSuspicious"
const precedingIPAccess = "precedingIpAccess"
const sebsequentIPAccess = "subsequentIpAccess"

// Radius of earth in miles
const earthRadius = 3949.9026

var db *sql.DB

// Config holds the configuration data for this database
type Config struct {
	Filename string `default:"./data/superman.db"`
}

// IPData contains the data returned to the user for a single IP - three of these in a return struct
type IPData struct {
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Radius    int     `json:"radius"`
	IP        string  `json:"ip,omitempty"`
	Timestamp int     `json:"timestamp,omitempty"`
	Speed     float64 `json:"speed,omitempty"`
}

// LocationRow contains the summation of data stored in a single row of the DB
type LocationRow struct {
	ID       string
	Username string
	IPData   IPData
}

// Setup sets up the DB connection
func Setup(conf Config) {
	fmt.Println("Haversine DB connection setup")

	var err error
	//TODO: review the Open() function in the library: does it work like databases/sql?
	db, err = sql.Open("sqlite3", conf.Filename)
	if err != nil {
		fmt.Printf("Haversine database initialization failed with error: %+v\n", err)
		return
	}
}

// GenerateResponse holds our central business logic
func GenerateResponse(c *gin.Context) {
	fmt.Println("Business logic")
	//Ping DB to see if it works
	err := db.Ping()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}

	//TODO: formalize this struct
	toreturn := map[string]interface{}{}

	//Fetch data from context and make sure it fits the data
	inputInterface, exists := c.Get(sanitize.InputContextKey)
	if !exists {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]interface{}{"error": "input context key doesn't exist"})
		return
	}
	input, ok := inputInterface.(sanitize.InputData)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]interface{}{"error": "bad input structure type"})
		return
	}
	locationInterface, exists := c.Get(geolocation.LocationContextKey)
	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{"error": "location context key doesn't exist"})
		return
	}
	location, ok := locationInterface.(geolocation.Loc)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]interface{}{"error": "bad location structure type"})
		return
	}

	currRow := LocationRow{
		ID:       input.UUID,
		Username: input.Username,
		IPData: IPData{
			Lon:       location.Longitude,
			Lat:       location.Latitude,
			Radius:    location.AccuracyRadius,
			IP:        input.IP,
			Timestamp: input.Timestamp,
		},
	}

	toreturn[currentGeo] = IPData{
		Lon:    location.Longitude,
		Lat:    location.Latitude,
		Radius: location.AccuracyRadius,
	}

	fmt.Printf("%+v %+v\n", toreturn, currRow)

	//Retrieve previous login data from DB based on timestamp of current record
	prevRow := LocationRow{}
	prevRowExists := false
	prevRowQuery, err := db.Query("SELECT id, username, timestamp, lat, lon, radius, ip FROM locations WHERE username = ? AND timestamp > ? ORDER BY timestamp DESC LIMIT 1", currRow.Username, currRow.IPData.Timestamp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	defer prevRowQuery.Close()
	for prevRowQuery.Next() {
		prevRowExists = true
		err := prevRowQuery.Scan(&prevRow.ID, &prevRow.Username, &prevRow.IPData.Timestamp, &prevRow.IPData.Lat, &prevRow.IPData.Lon, &prevRow.IPData.Radius, &prevRow.IPData.IP)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
			return
		}
		fmt.Printf("%+v\n", prevRow)
	}
	err = prevRowQuery.Err()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	if prevRowExists {
		fmt.Println("found previous row")
		//compute haversine, put speed and suspiciousness into the return struct
		var suspicious bool
		prevRow.IPData.Speed, suspicious = computeHaversine(prevRow.IPData.Lat, prevRow.IPData.Lon, currRow.IPData.Lat, currRow.IPData.Lon, currRow.IPData.Timestamp, prevRow.IPData.Timestamp)
		toreturn[precedingIPAccess] = prevRow.IPData
		toreturn[travelToCurrentGeoSuspicious] = suspicious
	}

	//Aaaand retrieve next login data from DB based on timestamp of current record
	nextRow := LocationRow{}
	nextRowExists := false
	nextRowQuery, err := db.Query("SELECT id, username, timestamp, lat, lon, radius, ip FROM locations WHERE username = ? AND timestamp < ? ORDER BY timestamp ASC LIMIT 1", currRow.Username, currRow.IPData.Timestamp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	defer nextRowQuery.Close()
	for nextRowQuery.Next() {
		nextRowExists = true
		err := nextRowQuery.Scan(&nextRow.ID, &nextRow.Username, &nextRow.IPData.Timestamp, &nextRow.IPData.Lat, &nextRow.IPData.Lon, &nextRow.IPData.Radius, &nextRow.IPData.IP)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
			return
		}
		fmt.Printf("%+v\n", nextRow)
	}
	err = nextRowQuery.Err()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	if nextRowExists {
		fmt.Println("found next row")
		//compute haversine, put speed and suspiciousness into the return struct
		var suspicious bool
		nextRow.IPData.Speed, suspicious = computeHaversine(currRow.IPData.Lat, currRow.IPData.Lon, nextRow.IPData.Lat, nextRow.IPData.Lon, currRow.IPData.Timestamp, nextRow.IPData.Timestamp)
		toreturn[sebsequentIPAccess] = nextRow.IPData
		toreturn[travelFromCurrentGeoSuspicious] = suspicious
	}

	//Insert the model into the DB - defer() this? what to do with multiple timestamps that are the same, due to unique constraint on table?
	_, err = db.Exec("INSERT INTO locations (id, username, timestamp, lat, lon, radius, ip) VALUES (?,?,?,?,?,?,?)", currRow.ID, currRow.Username, currRow.IPData.Timestamp, currRow.IPData.Lat, currRow.IPData.Lon, currRow.IPData.Radius, currRow.IPData.IP)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}

	//finally, return what we've got
	c.JSON(http.StatusOK, toreturn)
}

func computeHaversine(lat1, lon1, lat2, lon2 float64, currTime, prevTime int) (speed float64, suspcious bool) {
	lat1 = lat1 * math.Pi / 180
	lat2 = lat2 * math.Pi / 180
	lon1 = lon1 * math.Pi / 180
	lon2 = lon2 * math.Pi / 180

	cosinesPart := math.Cos(lat1) * math.Cos(lat2) * math.Pow(math.Sin((lon2-lon1)/2), 2)
	sinPart := math.Pow((math.Sin(lat2-lat1) / 2), 2)
	arcsinPart := math.Asin(math.Sqrt(sinPart + cosinesPart))
	distance := 2 * earthRadius * arcsinPart

	timeInHours := (float64(currTime) - float64(prevTime)) / 3600
	speed = distance / timeInHours
	if speed > 500 {
		return speed, true
	}
	return speed, false
}
