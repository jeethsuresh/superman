package geolocation

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeethsuresh/superman/pkg/sanitize"
	"github.com/oschwald/geoip2-golang"
)

// LocationContextKey is used in other packages to find the location struct within the context
const LocationContextKey = "LOCATION"

var db *geoip2.Reader

// Config contains any information required for this middleware to work
type Config struct {
	Filename string `default:"./data/GeoLite2-City.mmdb"`
}

// Loc defines a structure for the important information extracted from the GeoIP database
type Loc struct {
	Latitude       float64
	Longitude      float64
	AccuracyRadius uint16
}

// Setup is run once, to set up the database connection from the config env var
func Setup(conf Config) {
	fmt.Println("Geolocation middleware setup")

	var err error
	//TODO: review the Open() function in the library: does it work like databases/sql?
	db, err = geoip2.Open(conf.Filename)
	if err != nil {
		fmt.Println("database initialization failed")
		return
	}
}

// Middleware is the handler func that slots into gin's middleware chain
func Middleware(c *gin.Context) {
	fmt.Println("Geolocation Middleware")

	inputInterface, exists := c.Get(sanitize.InputContextKey)
	if !exists {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]interface{}{"error": "context key doesn't exist"})
		return
	}
	input, ok := inputInterface.(sanitize.InputData)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]interface{}{"error": "bad input structure type"})
		return
	}

	ip := net.ParseIP(input.IP)
	if ip == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid IP"})
		return
	}
	record, err := db.City(ip)
	if err != nil {
		fmt.Printf("[ERR] %+v\n", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	Location := Loc{
		Latitude:       record.Location.Latitude,
		Longitude:      record.Location.Longitude,
		AccuracyRadius: record.Location.AccuracyRadius,
	}
	c.Set(LocationContextKey, Location)
	fmt.Printf("Geolocation: %+v\n", Location)
	c.Next()
}
