package calendar

import (
	"strings"
	"sync"
	"time"
)

var (
	lock           sync.Mutex
	installedClock Clock = &WallClock{}

	namedLayouts = map[string]string{
		"ANSIC":       time.ANSIC,
		"UnixDate":    time.UnixDate,
		"RubyDate":    time.RubyDate,
		"RFC822":      time.RFC822,
		"RFC822Z":     time.RFC822Z,
		"RFC850":      time.RFC850,
		"RFC1123":     time.RFC1123,
		"RFC1123Z":    time.RFC1123Z,
		"RFC3339":     time.RFC3339,
		"RFC3339Nano": time.RFC3339Nano,
		"Kitchen":     time.Kitchen,
		"Stamp":       time.Stamp,
		"StampMilli":  time.StampMilli,
		"StampMicro":  time.StampMicro,
		"StampNano":   time.StampNano,
		"DateTime":    time.DateTime,
		"DateOnly":    time.DateOnly,
		"TimeOnly":    time.TimeOnly,
	}
)

type Clock interface {
	// Returns the current local time
	Now() time.Time

	// Returns the current UTC time
	UtcNow() time.Time
}

// A clock that returns the current wallclock time
type WallClock struct{}

func (clock *WallClock) Now() time.Time {
	return time.Now()
}

func (clock *WallClock) UtcNow() time.Time {
	return time.Now().UTC()
}

// Registers a new system wide clock
func RegisterClock(clock Clock) Clock {
	lock.Lock()
	defer lock.Unlock()

	var oldClock = installedClock
	installedClock = clock

	return oldClock
}

// Returns the current local time from the system wide clock
func Now() time.Time {
	lock.Lock()
	clock := installedClock
	lock.Unlock()

	return clock.Now()
}

// Returns the current UTC time from the system wide clock
func UtcNow() time.Time {
	lock.Lock()
	clock := installedClock
	lock.Unlock()

	return clock.UtcNow()
}

// Converts more user friendly date and time formats to Go format
func FormatConverter(layout string) string {
	if format, ok := namedLayouts[layout]; ok {
		return format
	}

	layout = strings.ReplaceAll(layout, "YYYY", "2006")
	layout = strings.ReplaceAll(layout, "YY", "06")

	layout = strings.ReplaceAll(layout, "MMMM", "January")
	layout = strings.ReplaceAll(layout, "MMM", "Jan")
	layout = strings.ReplaceAll(layout, "MM", "01")
	layout = strings.ReplaceAll(layout, "M", "1")

	layout = strings.ReplaceAll(layout, "DDDD", "Monday")
	layout = strings.ReplaceAll(layout, "DDD", "Mon")
	layout = strings.ReplaceAll(layout, "DD", "02")
	layout = strings.ReplaceAll(layout, "D", "2")

	layout = strings.ReplaceAll(layout, "HH", "15")
	layout = strings.ReplaceAll(layout, "hh", "03")
	layout = strings.ReplaceAll(layout, "h", "3")

	layout = strings.ReplaceAll(layout, "mm", "04")
	layout = strings.ReplaceAll(layout, "m", "4")

	layout = strings.ReplaceAll(layout, "ss", "05")
	layout = strings.ReplaceAll(layout, "s", "5")

	layout = strings.ReplaceAll(layout, "fffffffff", "000000000")
	layout = strings.ReplaceAll(layout, "ffffff", "000000")
	layout = strings.ReplaceAll(layout, "fff", "000")

	return layout
}
