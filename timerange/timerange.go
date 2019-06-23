package timerange

import (
	"errors"
	"strings"
	"time"
)

type TimeRange struct {
	start time.Time
	end time.Time
}

// Create a new TimeRange by immediately parsing the incoming string
func NewTimeRange(inputTimeRange string) (*TimeRange, error) {
	value := TimeRange{}

	err := value.Parse(inputTimeRange)

	if err != nil {
		return nil, err
	}

	return &value, nil
}

/**
 * Get two GO times as a range from a string
 * Note both the start and end times will be set to the current day,
 *  so it is possible for end to actually be before start.. which is OK as we're not really interested
 *  in the Date, and it's that the TimeRange will be comparing dates against ONLY the current day
 */
func (t *TimeRange) Parse(inputTimeRange string) error {
	timeRangeSplit := strings.Split(inputTimeRange, "-")

	if len(timeRangeSplit) != 2 {
		return errors.New("timeRange does not have exactly two parts")
	}

	startTime, err := time.Parse("15:04", timeRangeSplit[0])
	if err != nil {
		// Bubble the error up to the caller
		return err
	}
	endTime, err := time.Parse("15:04", timeRangeSplit[1])
	if err != nil {
		// Bubble the error up to the caller
		return err
	}

	now := time.Now()
	t.start = time.Date(now.Year(), now.Month(), now.Day(), startTime.Hour(), startTime.Minute(), 0, 0, time.Local)
	t.end = time.Date(now.Year(), now.Month(), now.Day(), endTime.Hour(), endTime.Minute(), 0, 0, time.Local)

	// No error.. yay!
	return nil
}

// Is a given time within the Time Range and on this day
func (t *TimeRange) HasTime(compareTime time.Time) bool {
	// First check the date is even on the same day
	if compareTime.Year() != t.start.Year() ||
		compareTime.Month() != t.start.Month() ||
		compareTime.Day() != t.start.Day() {
		return false
	}

	if compareTime.Equal(t.start) {
		return true
	}

	// First case, the compare time is between two times on the same day (e.g. 10am -> 11am)
	if t.end.After(t.start) {
		return compareTime.After(t.start) && compareTime.Before(t.end)
	}

	// Second case, the end time is before the start time (e.g. on the next day - 23pm -> 1am)
	if t.end.Before(t.start) {
		return compareTime.After(t.start) || compareTime.Before(t.end)
	}

	return false
}

func (t *TimeRange) Print() string {
	return t.start.String() + " - " + t.end.String()
}