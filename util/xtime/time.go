package xtime

import (
	"time"
)

// Time time
type Time struct {
	time.Time
}

// Now returns current time
func Now() *Time {
	return &Time{
		Time: time.Now(),
	}
}

// Unix returns time converted from timestamp
func Unix(sec, nsec int64) *Time {
	return &Time{
		Time: time.Unix(sec, nsec),
	}
}

// Today returns begin time of today
func Today() *Time {
	return Now().BeginOfDay()
}

// BeginOfYear BeginOfYear beginning of year
func (t *Time) BeginOfYear() *Time {
	y, _, _ := t.Date()
	return &Time{time.Date(y, time.January, 1, 0, 0, 0, 0, t.Location())}
}

// EndOfYear end of year
func (t *Time) EndOfYear() *Time {
	return &Time{t.BeginOfYear().AddDate(1, 0, 0).Add(-time.Nanosecond)}
}

// BeginOfMonth begin of month
func (t *Time) BeginOfMonth() *Time {
	y, m, _ := t.Date()
	return &Time{time.Date(y, m, 1, 0, 0, 0, 0, t.Location())}
}

// EndOfMonth end of month
func (t *Time) EndOfMonth() *Time {
	return &Time{t.BeginOfMonth().AddDate(0, 1, 0).Add(-time.Nanosecond)}
}

// BeginOfWeek the day of begin of week,
// NOTE: week begin from Sunday
func (t *Time) BeginOfWeek() *Time {
	y, m, d := t.AddDate(0, 0, 0-int(t.BeginOfDay().Weekday())).Date()
	return &Time{time.Date(y, m, d, 0, 0, 0, 0, t.Location())}
}

// EndOfWeek end of week
// NOTE: week end with Saturday
func (t *Time) EndOfWeek() *Time {
	y, m, d := t.BeginOfWeek().AddDate(0, 0, 7).Add(-time.Nanosecond).Date()
	return &Time{time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())}
}

// BeginOfDay returns zero point of time's day
func (t *Time) BeginOfDay() *Time {
	y, m, d := t.Date()
	return &Time{time.Date(y, m, d, 0, 0, 0, 0, t.Location())}
}

// EndOfDay returns last point of time's day
func (t *Time) EndOfDay() *Time {
	y, m, d := t.Date()
	return &Time{time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())}
}

// BeginOfHour returns zero point of time's day
func (t *Time) BeginOfHour() *Time {
	y, m, d := t.Date()
	return &Time{time.Date(y, m, d, t.Hour(), 0, 0, 0, t.Location())}
}

// EndOfHour returns last point of time's day
func (t *Time) EndOfHour() *Time {
	y, m, d := t.Date()
	return &Time{time.Date(y, m, d, t.Hour(), 59, 59, int(time.Second-time.Nanosecond), t.Location())}
}

// BeginOfMinute returns zero point of time's day
func (t *Time) BeginOfMinute() *Time {
	y, m, d := t.Date()
	return &Time{time.Date(y, m, d, t.Hour(), t.Minute(), 0, 0, t.Location())}
}

// EndOfMinute returns last point of time's day
func (t *Time) EndOfMinute() *Time {
	y, m, d := t.Date()
	return &Time{time.Date(y, m, d, t.Hour(), t.Minute(), 59, int(time.Second-time.Nanosecond), t.Location())}
}

var TS TimeFormat = "2006-01-02 15:04:05"

type TimeFormat string

func (ts TimeFormat) Format(t time.Time) string {
	return t.Format(string(ts))
}
