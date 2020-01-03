package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/jeethsuresh/superman/pkg/geolocation"
	"github.com/jeethsuresh/superman/pkg/haversine"
	"github.com/jeethsuresh/superman/pkg/sanitize"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	fmt.Println("Beginning Superman detection program")

	r := gin.New()

	r.Use(sanitize.Middleware)

	var geolocationConfig geolocation.Config
	err := envconfig.Process("superman_geo", &geolocationConfig)
	if err != nil {
		panic("Failure to process Geolocation Config")
	}
	geolocation.Setup(geolocationConfig)
	r.Use(geolocation.Middleware)

	var haversineConfig haversine.Config
	err = envconfig.Process("superman_geo", &haversineConfig)
	if err != nil {
		panic("Failure to process Haversine Location DB Config")
	}
	haversine.Setup(haversineConfig)
	r.POST("/v1", haversine.GenerateResponse)

	fmt.Println("==========================")
	_ = r.Run(":5000")
}
