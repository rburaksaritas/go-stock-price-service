package utils

import (
	"fmt"
	"time"
)

var timeZones = map[string]string{
	"tr": "Europe/Istanbul",
	"us": "America/New_York",
	"uk": "Europe/London",
	"jp": "Asia/Tokyo",
	"in": "Asia/Kolkata",
}

func LoadLocation(timeZone string) (*time.Location, error) {
	if tz, ok := timeZones[timeZone]; ok {
		return time.LoadLocation(tz)
	}
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return nil, fmt.Errorf("could not load location: %v", err)
	}
	return loc, nil
}

func Int64ToReadableTimestamp(timestamp int64, timezone string) (*string, error) {
	loc, err := LoadLocation(timezone)
	if err != nil {
		return nil, err
	}

	readableTimestamp := time.Unix(timestamp, 0).In(loc).Format(time.RFC3339)
	return &readableTimestamp, nil
}

func TimeToReadableTimestamp(timeObj time.Time, timezone string) (*string, error) {
	loc, err := LoadLocation(timezone)
	if err != nil {
		return nil, err
	}

	readableTimestamp := timeObj.In(loc).Format(time.RFC3339)
	return &readableTimestamp, nil
}
