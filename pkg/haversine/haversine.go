package haversine

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeethsuresh/superman/pkg/geolocation"
	"github.com/jeethsuresh/superman/pkg/sanitize"

	//Well-known SQLite3 driver
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Config holds the configuration data for this database
type Config struct {
	Filename string `default:"./data/superman.db"`
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
	}

	//Fetch data from context
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

	fmt.Printf("Input: %+v Location: %+v\n", input, location)

	//Retrieve previous/next login data from DB based on timestamp of current record

	//If either timestamp returned is the same, return an invalid request if IPs are different or a not-suspicious if they're the same

	//Compute haversine distance for previous (if exists) and next (if exists) location
	//  absolute values?
	//  moving backwards?
	//

	//If necessary, update the previous and next entries' speed in the DB

	//Insert the model into the DB - defer() this?

	var id, username string
	rows, err := db.Query("select id, username from locations where id = ?", 1)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &username)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}
		fmt.Printf("%+v %+v\n", id, username)
	}
	err = rows.Err()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}
}
