package utils

import (
	"fmt"
	"time"
)

const (
	dateTime = "2006-01-02 15:04"
	dateOnly = "2006-01-02"
)

// LocalDateTime converts the given date time to the given timezone
func LocalDateTime(dateTime time.Time, timezone string) time.Time {
	loc, _ := time.LoadLocation(timezone)
	return dateTime.In(loc)
}

// TimeStringToCurrentDateTime parses the given time string and merges it with the current local date
func TimeStringToCurrentDateTime(timeString string) (time.Time, error) {
	currentDate := time.Now().Local().Format(dateOnly)
	parsedDateTime, err := time.Parse(dateTime, currentDate+" "+timeString)
	if err != nil {
		return time.Time{}, fmt.Errorf("Failed to parse time string to date time, %w", err)
	}
	return parsedDateTime, nil
}

// ParseDateTime parses the given date string to a time.Time object according to the RFC3339 format
func ParseDateTime(dateString string) (time.Time, error) {
	return time.Parse(time.RFC3339, dateString)
}

// DateTimeToString converts the given date time to a string according to the RFC3339 format
func DateTimeToString(dateTime time.Time) string {
	return dateTime.Format(time.RFC3339)
}
