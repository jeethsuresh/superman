package haversine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeHaversineValidAustinToHouston(t *testing.T) {
	//Valid coordinates for Austin and Houston, plus a large radius
	lat1 := 30.26715
	lon1 := -97.74306
	lat2 := 29.76328
	lon2 := -95.36327
	time1 := 1576396800
	time2 := 1576411200
	rad1 := 20
	rad2 := 20

	speed, suspicious := computeHaversine(lat1, lon1, lat2, lon2, time1, time2, rad1, rad2)
	assert.True(t, (speed < 500))
	assert.False(t, suspicious)
}
func TestComputeHaversineValidHoustonToAustin(t *testing.T) {
	// Same coordinates as previous test, but with the times reversed
	lat1 := 30.26715
	lon1 := -97.74306
	lat2 := 29.76328
	lon2 := -95.36327
	time2 := 1576396800
	time1 := 1576411200
	rad1 := 20
	rad2 := 20

	_, suspicious := computeHaversine(lat1, lon1, lat2, lon2, time1, time2, rad1, rad2)
	assert.False(t, suspicious)
}
func TestComputeHaversineInvalidAustinToHouston(t *testing.T) {
	// Same coordinates as previous test, but with bad times
	lat1 := 30.26715
	lon1 := -97.74306
	lat2 := 29.76328
	lon2 := -95.36327
	time1 := 1576396800
	time2 := 1576396801
	rad1 := 20
	rad2 := 20

	speed, suspicious := computeHaversine(lat1, lon1, lat2, lon2, time1, time2, rad1, rad2)
	assert.False(t, (speed < 500))
	assert.True(t, suspicious)
}

func TestComputeHaversineBadDistance(t *testing.T) {
	//Valid coordinates for Austin and Houston, plus a comically large radius for Austin IPs
	lat1 := 30.26715
	lon1 := -97.74306
	lat2 := 29.76328
	lon2 := -95.36327
	time1 := 1576396800
	time2 := 1576411200
	rad1 := 200
	rad2 := 20

	//This test is an edge case; I'm going to say it's not suspicious, because there's a chance that they're in the overlapping radius and their IPs have switched
	speed, suspicious := computeHaversine(lat1, lon1, lat2, lon2, time1, time2, rad1, rad2)
	assert.True(t, (speed < 500))
	assert.False(t, suspicious)
}

func TestComputeHaversineBadDistanceShortTime(t *testing.T) {
	//Valid coordinates for Austin and Houston, plus a comically large radius for Austin IPs
	lat1 := 30.26715
	lon1 := -97.74306
	lat2 := 29.76328
	lon2 := -95.36327
	time1 := 1576396800
	time2 := 1576396801
	rad1 := 200
	rad2 := 20

	//This test expands on the idea that two simultaneous logins for the same user may come from two different IP blocks whose radii overlap, which still shouldn't be flagged as suspicious (your phone vs your laptop)
	speed, suspicious := computeHaversine(lat1, lon1, lat2, lon2, time1, time2, rad1, rad2)
	assert.True(t, (speed < 500))
	assert.False(t, suspicious)
}

func TestComputeHaversine500mphAustinToHouston(t *testing.T) {
	//Valid coordinates for Austin and Houston, plus a large radius
	lat1 := 30.26715
	lon1 := -97.74306
	lat2 := 29.76328
	lon2 := -95.36327
	time1 := 1576396800
	time2 := 1576397392
	rad1 := 20
	rad2 := 20

	//Speed just above 500mph should still trigger
	speed, suspicious := computeHaversine(lat1, lon1, lat2, lon2, time1, time2, rad1, rad2)
	assert.False(t, (speed < 500))
	assert.True(t, suspicious)
}

func TestComputeHaversine499mphAustinToHouston(t *testing.T) {
	//Valid coordinates for Austin and Houston, plus a large radius
	lat1 := 30.26715
	lon1 := -97.74306
	lat2 := 29.76328
	lon2 := -95.36327
	time1 := 1576396800
	time2 := 1576397393
	rad1 := 20
	rad2 := 20

	//Speed just below 500mph should not trigger
	speed, suspicious := computeHaversine(lat1, lon1, lat2, lon2, time1, time2, rad1, rad2)
	assert.True(t, (speed < 500))
	assert.False(t, suspicious)
}
