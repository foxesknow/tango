package calendar

import (
	"sync"
	"time"
)

var (
	lock           sync.Mutex
	installedClock Clock = &WallClock{}
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
