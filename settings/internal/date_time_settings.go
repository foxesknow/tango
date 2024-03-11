package providers

import (
	"strings"
	"time"

	"github.com/foxesknow/tango/calendar"
)

type formatableTime struct {
	t time.Time
}

func (f *formatableTime) String() string {
	return f.t.String()
}

func (f *formatableTime) Format(layout string) string {
	convertedLayout := calendar.FormatConverter(layout)
	return f.t.Format(convertedLayout)
}

type DateTimeSettings struct {
}

func (self *DateTimeSettings) GetSetting(name string) (value any, found bool) {
	switch strings.ToLower(name) {
	case "now":
		return &formatableTime{t: calendar.Now()}, true
	case "utcnow":
		return &formatableTime{t: calendar.UtcNow()}, true
	default:
		return "", false
	}
}
