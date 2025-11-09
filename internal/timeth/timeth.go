package timeth

import (
	"time"

	log "github.com/sirupsen/logrus"
)

var timezoneLocation *time.Location

func init() {
	var err error
	timezoneLocation, err = time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Warn("Failed to load Asia/Bangkok timezone, using UTC")
		timezoneLocation = time.UTC
	}
}

func Now() time.Time {
	return time.Now().In(timezoneLocation)
}

func LoadLocation() *time.Location {
	return timezoneLocation
}
