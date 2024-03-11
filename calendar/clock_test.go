package calendar

import (
	"testing"
	"time"
)

func Test_Now_Is_Valid(t *testing.T) {
	now := Now()

	var zero time.Time
	if now == zero {
		t.Fatal("invalid time")
	}
}
