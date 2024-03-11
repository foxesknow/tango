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

func Test_UtcNow_Is_Valid(t *testing.T) {
	now := UtcNow()

	var zero time.Time
	if now == zero {
		t.Fatal("invalid time")
	}
}

func Test_RegisterClock(t *testing.T) {
	old := RegisterClock(&WallClock{})
	defer RegisterClock(old)

	now := UtcNow()

	var zero time.Time
	if now == zero {
		t.Fatal("invalid time")
	}
}

func Test_FormatConverter_Years(t *testing.T) {
	now := time.Date(2024, time.March, 6, 7, 8, 9, 0, time.Local)

	if now.Format(FormatConverter("YYYY")) != "2024" {
		t.Fatal("format failed")
	}

	if now.Format(FormatConverter("YY")) != "24" {
		t.Fatal("format failed")
	}
}

func Test_FormatConverter_Month(t *testing.T) {
	now := time.Date(2024, time.March, 6, 7, 8, 9, 0, time.Local)

	if now.Format(FormatConverter("MMMM")) != "March" {
		t.Fatal("format failed")
	}

	if now.Format(FormatConverter("MMM")) != "Mar" {
		t.Fatal("format failed")
	}

	if now.Format(FormatConverter("MM")) != "03" {
		t.Fatal("format failed")
	}

	if now.Format(FormatConverter("M")) != "3" {
		t.Fatal("format failed")
	}
}

func Test_FormatConverter_Day(t *testing.T) {
	now := time.Date(2024, time.March, 6, 7, 8, 9, 0, time.Local)

	if now.Format(FormatConverter("DDDD")) != "Wednesday" {
		t.Fatal("format failed")
	}

	if now.Format(FormatConverter("DDD")) != "Wed" {
		t.Fatal("format failed")
	}

	if now.Format(FormatConverter("DD")) != "06" {
		t.Fatal("format failed")
	}

	if now.Format(FormatConverter("D")) != "6" {
		t.Fatal("format failed")
	}
}

func Test_FormatConverter_24_Hour(t *testing.T) {
	now := time.Date(2024, time.March, 6, 14, 8, 9, 0, time.Local)

	if now.Format(FormatConverter("HH")) != "14" {
		t.Fatal("format failed")
	}
}

func Test_FormatConverter_Hour(t *testing.T) {
	now := time.Date(2024, time.March, 6, 7, 8, 9, 0, time.Local)

	if now.Format(FormatConverter("HH")) != "07" {
		t.Fatal("format failed")
	}

	if now.Format(FormatConverter("hh")) != "07" {
		t.Fatal("format failed")
	}

	if now.Format(FormatConverter("h")) != "7" {
		t.Fatal("format failed")
	}
}

func Test_FormatConverter_Minute(t *testing.T) {
	now := time.Date(2024, time.March, 6, 7, 8, 9, 0, time.Local)

	if now.Format(FormatConverter("mm")) != "08" {
		t.Fatal("format failed")
	}

	if now.Format(FormatConverter("m")) != "8" {
		t.Fatal("format failed")
	}
}

func Test_FormatConverter_Second(t *testing.T) {
	now := time.Date(2024, time.March, 6, 7, 8, 9, 0, time.Local)

	if now.Format(FormatConverter("ss")) != "09" {
		t.Fatal("format failed")
	}

	if now.Format(FormatConverter("s")) != "9" {
		t.Fatal("format failed")
	}
}

func Test_FormatConverter_StandardLayout_ANSIC(t *testing.T) {
	now := time.Date(2024, time.March, 6, 7, 8, 9, 0, time.Local)

	if now.Format(FormatConverter("ANSIC")) != "Wed Mar  6 07:08:09 2024" {
		t.Fatal("format failed")
	}
}

func Test_FormatConverter_StandardLayout_Kitchen(t *testing.T) {
	now := time.Date(2024, time.March, 6, 7, 8, 9, 0, time.Local)

	if now.Format(FormatConverter("Kitchen")) != "7:08AM" {
		t.Fatal("format failed")
	}
}
