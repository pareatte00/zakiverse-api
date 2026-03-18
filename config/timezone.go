package config

import (
	"log"
	"time"
)

func setGlobalTimezone(tz string) error {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Fatalf("Failed to load location timezone config : %s", err)
	}

	time.Local = loc

	return nil
}
