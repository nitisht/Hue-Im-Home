package timerange

import (
	"testing"
	"time"
)

func TestParseValidTimeRange(t *testing.T) {
	timeRange, err := NewTimeRange("23:12-01:30")

	if err != nil {
		t.Fatal(err.Error())
	}

	if timeRange.start.Hour() != 23 || timeRange.start.Minute() != 12 {
		t.Fatal("Failed to parse Start Time")
	}

	if timeRange.end.Hour() != 1 || timeRange.end.Minute() != 30 {
		t.Fatal("Failed to parse End Time")
	}
}

func TestParseEmptyTimeRange(t *testing.T) {
	_, err := NewTimeRange("")

	if err == nil {
		t.Fatal("Failed to handle empty time range")
	}
}

func TestTimeWithinTimeRangeSameDay(t *testing.T) {
	// First test simple time comparison
	timeRange, _ := NewTimeRange("10:00-11:00")
	newTime := getDateForTime(10, 1)

	if !timeRange.HasTime(newTime) {
		t.Fail()
	}
}

func TestTimeWithinTimeRangeNextDay(t *testing.T) {
	// Now test end time before start time (e.g. overnight)
	timeRange, _ := NewTimeRange("23:00-01:00")
	newTime := getDateForTime(23, 30)
	if !timeRange.HasTime(newTime) {
		t.Fail()
	}
}

func TestTimeWithinTimeRangeNextDayAlt(t *testing.T) {
	// Now test end time before start time (e.g. overnight)
	timeRange, _ := NewTimeRange("23:00-02:00")
	newTime := getDateForTime(1, 30)
	if !timeRange.HasTime(newTime) {
		t.Fail()
	}
}

func TestTimeOutsideTimeRangeSameDay(t *testing.T) {
	// First test simple time comparison
	timeRange, _ := NewTimeRange("10:00-11:00")
	newTime := getDateForTime(11, 00)
	if timeRange.HasTime(newTime) {
		t.Fail()
	}
}

func TestTimeOutsideTimeRangeNextDay(t *testing.T) {
	// Now test end time before start time (e.g. overnight)
	timeRange, _ := NewTimeRange("23:00-01:00")
	newTime := getDateForTime(22, 59)
	if timeRange.HasTime(newTime) {
		t.Fail()
	}
}

func TestTimeOutsideTimeRangeNextDayAlt(t *testing.T) {
	// Now test end time before start time (e.g. overnight)
	timeRange, _ := NewTimeRange("23:00-01:00")
	newTime := getDateForTime(1, 10)
	if timeRange.HasTime(newTime) {
		t.Fail()
	}
}

func getDateForTime(hour int, minute int) time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour, minute, 0, 0, time.Local)
}